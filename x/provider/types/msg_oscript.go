package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	msgSetOScript  = "set_oscript"
	msgEditOScript = "edit_oscript"
)

// MsgCreateOracleScript defines a CreateOracleScript message
type MsgCreateOracleScript struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
}

// NewMsgCreateOracleScript is a constructor function for MsgCreateOracleScript
func NewMsgCreateOracleScript(name string, code []byte, owner sdk.AccAddress, des string) MsgCreateOracleScript {
	return MsgCreateOracleScript{
		Name:        name,
		Code:        code,
		Owner:       owner,
		Description: des,
	}
}

// Route should return the name of the module
func (msg MsgCreateOracleScript) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateOracleScript) Type() string { return msgSetOScript }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateOracleScript) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Name) == 0 || len(msg.Code) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditOracleScript defines a message for editing a oScript in the store
type MsgEditOracleScript struct {
	OldName     string         `json:"old_name"`
	NewName     string         `json:"new_name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
}

// NewMsgEditOracleScript is a constructor function for MsgEditOracleScript
func NewMsgEditOracleScript(oldName string, newName string, code []byte, owner sdk.AccAddress, des string) MsgEditOracleScript {
	return MsgEditOracleScript{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Code:        code,
		Owner:       owner,
	}
}

// Route should return the name of the module
func (msg MsgEditOracleScript) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditOracleScript) Type() string { return msgEditOScript }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditOracleScript) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Code) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgEditOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}