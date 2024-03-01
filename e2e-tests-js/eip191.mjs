import { makeSignDoc, serializeSignDoc,coins, encodeSecp256k1Signature, encodeSecp256k1Pubkey } from "@cosmjs/amino";
import { StargateClient, calculateFee, createDefaultAminoConverters, AminoTypes, defaultRegistryTypes as defaultStargateTypes} from "@cosmjs/stargate";
import { createWasmAminoConverters, wasmTypes} from "@cosmjs/cosmwasm-stargate";
import "dotenv/config";
import { ethers } from "ethers";
import { keccak256, ripemd160, sha256 } from '@cosmjs/crypto';
import { fromHex, toUtf8, toHex, fromBase64 } from '@cosmjs/encoding';
import { TxRaw } from 'cosmjs-types/cosmos/tx/v1beta1/tx.js';
import { SignMode } from 'cosmjs-types/cosmos/tx/signing/v1beta1/signing.js';
import {recoverPublicKey} from '@noble/secp256k1';
import bech32 from 'bech32';
import { Registry, makeAuthInfoBytes , encodePubkey} from "@cosmjs/proto-signing";
import { Int53 } from '@cosmjs/math';

async function startMultinodeLocal() {
  await executeSpawn(path.join(scriptPath, "multinode-local-testnet.sh"));
}

async function cleanNetwork() {
  await executeSpawn(path.join(scriptPath, "clean-multinode-local-testnet.sh"));
}

export function pubkeyToBechAddress(pubkey, prefix= 'orai'){
  return bech32.encode(prefix, bech32.toWords(ripemd160(sha256(pubkey))));
}

export function getPubkeyFromEthSignature(rawMsg, sigResult){

  // On ETHland pubkeys are recovered from signatures, so we're going to:
  // 1. sign something
  // 2. recover the pubkey from the signature
  // 3. derive a secret address from the the pubkey

  // strip leading 0x and extract recovery id
  const sig = fromHex(sigResult.slice(2, -2));
  let recoveryId = parseInt(sigResult.slice(-2), 16) - 27;

  // When a Ledger is used, this value doesn't need to be adjusted
  if (recoveryId < 0) {
    recoveryId += 27;
  }
  const eip191MessagePrefix = toUtf8('\x19Ethereum Signed Message:\n');
  const rawMsgLength = toUtf8(String(rawMsg.length));

  const publicKey = recoverPublicKey(keccak256(new Uint8Array([...eip191MessagePrefix, ...rawMsgLength, ...rawMsg])), sig, recoveryId, true);

  return publicKey;
}
async function eip191SendTokens(){
  const amount = 1000000;
  const target_address = "orai1cnza7u4g9nwl5algvjfzwdlry2gk8andwgh4q8";

  const walletEthers = new ethers.Wallet(process.env.PRIVATE_KEY);
  const etherPublicKey = ethers.utils.computePublicKey(walletEthers.publicKey,true);
  const arrayPubKey = ethers.utils.arrayify(etherPublicKey);

  const cosmosAddress = pubkeyToBechAddress(arrayPubKey) 

  const signDoc = makeSignDoc(
     [
        {
          type: 'cosmos-sdk/MsgSend',
          value: {
            from_address: cosmosAddress,
            to_address: target_address,
            amount: coins(amount, 'orai')
          }
        }
      ],
      calculateFee(1000000, '0.002orai'),
      'testing',
      'memo',
      '10',
      '1'
  );
  // sign
  const rawMsg = serializeSignDoc(signDoc);
  // different from 'personal_sign' is that it sign the hex message directly
  const msgToSign = ethers.utils.arrayify(`0x${toHex(rawMsg)}`);
  const sigResult = await walletEthers.signMessage(msgToSign);
  const pubkey = getPubkeyFromEthSignature(rawMsg, sigResult);
  console.log("recoverPublicKey", pubkeyToBechAddress(pubkey))
  console.log("directFromWallet", cosmosAddress);
// strip leading 0x and trailing recovery id
  const sig = fromHex(sigResult.slice(2, -2));
  const signedMessage = {
    signed: signDoc,
    signature: encodeSecp256k1Signature(pubkey, sig)
  }

  const {signed, signature} = signedMessage;


  // build transaction
    const aminoTypes =   new AminoTypes({
        ...createDefaultAminoConverters(),
        ...createWasmAminoConverters()
      })
    const registry = new Registry([...defaultStargateTypes, ...wasmTypes])
    const signMode = SignMode.SIGN_MODE_EIP_191;
    const signedTxBody= {
      typeUrl: '/cosmos.tx.v1beta1.TxBody',
      value: {
        messages: signed.msgs.map((msg) => aminoTypes.fromAmino(msg)),
        memo: signed.memo
      }
    };
    const signedTxBodyBytes = registry.encode(signedTxBody);
    const signedGasLimit = Int53.fromString(signed.fee.gas).toNumber();
    const signedSequence = Int53.fromString(signed.sequence).toNumber();
    if (!pubkey) {
      throw new Error('Pubkey is required');
    }

    const any_pub_key = encodePubkey(encodeSecp256k1Pubkey(pubkey));

    const signedAuthInfoBytes = makeAuthInfoBytes([{ pubkey:any_pub_key, sequence: signedSequence }], signed.fee.amount, signedGasLimit, signed.fee.granter, signed.fee.payer, signMode);
    const txRaw = TxRaw.fromPartial({
      bodyBytes: signedTxBodyBytes,
      authInfoBytes: signedAuthInfoBytes,
      signatures: [fromBase64(signature.signature)]
    });
    const txBytes = TxRaw.encode(txRaw).finish();

  // broadcastTx
    const client = await StargateClient.connect("http://localhost:26657")
    const response = await client.broadcastTx(txBytes);
  console.log(response)
}

eip191SendTokens();
