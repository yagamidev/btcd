// Copyright (c) 2013 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcws

import (
	"encoding/json"
	"errors"
	"github.com/conformal/btcdb"
	"github.com/conformal/btcjson"
	"github.com/conformal/btcwire"
)

func init() {
	btcjson.RegisterCustomCmd("createencryptedwallet", parseCreateEncryptedWalletCmd)
	btcjson.RegisterCustomCmd("getbalances", parseGetBalancesCmd)
	btcjson.RegisterCustomCmd("getbestblock", parseGetBestBlockCmd)
	btcjson.RegisterCustomCmd("getcurrentnet", parseGetCurrentNetCmd)
	btcjson.RegisterCustomCmd("listalltransactions", parseListAllTransactionsCmd)
	btcjson.RegisterCustomCmd("notifynewtxs", parseNotifyNewTXsCmd)
	btcjson.RegisterCustomCmd("notifyspent", parseNotifySpentCmd)
	btcjson.RegisterCustomCmd("rescan", parseRescanCmd)
	btcjson.RegisterCustomCmd("walletislocked", parseWalletIsLockedCmd)
}

// GetCurrentNetCmd is a type handling custom marshaling and
// unmarshaling of getcurrentnet JSON websocket extension
// commands.
type GetCurrentNetCmd struct {
	id interface{}
}

// Enforce that GetCurrentNetCmd satisifies the btcjson.Cmd interface.
var _ btcjson.Cmd = &GetCurrentNetCmd{}

// NewGetCurrentNetCmd creates a new GetCurrentNetCmd.
func NewGetCurrentNetCmd(id interface{}) *GetCurrentNetCmd {
	return &GetCurrentNetCmd{id: id}
}

// parseGetCurrentNetCmd parses a RawCmd into a concrete type satisifying
// the btcjson.Cmd interface.  This is used when registering the custom
// command with the btcjson parser.
func parseGetCurrentNetCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 0 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	return NewGetCurrentNetCmd(r.Id), nil
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *GetCurrentNetCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *GetCurrentNetCmd) Method() string {
	return "getcurrentnet"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *GetCurrentNetCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "getcurrentnet",
		Id:      cmd.id,
	}
	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *GetCurrentNetCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseGetCurrentNetCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*GetCurrentNetCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// GetBestBlockCmd is a type handling custom marshaling and
// unmarshaling of getbestblock JSON websocket extension
// commands.
type GetBestBlockCmd struct {
	id interface{}
}

// Enforce that GetBestBlockCmd satisifies the btcjson.Cmd interface.
var _ btcjson.Cmd = &GetBestBlockCmd{}

// NewGetBestBlockCmd creates a new GetBestBlock.
func NewGetBestBlockCmd(id interface{}) *GetBestBlockCmd {
	return &GetBestBlockCmd{id: id}
}

// parseGetBestBlockCmd parses a RawCmd into a concrete type satisifying
// the btcjson.Cmd interface.  This is used when registering the custom
// command with the btcjson parser.
func parseGetBestBlockCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 0 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	return NewGetBestBlockCmd(r.Id), nil
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *GetBestBlockCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *GetBestBlockCmd) Method() string {
	return "getbestblock"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *GetBestBlockCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "getbestblock",
		Id:      cmd.id,
	}
	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *GetBestBlockCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseGetBestBlockCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*GetBestBlockCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// RescanCmd is a type handling custom marshaling and
// unmarshaling of rescan JSON websocket extension
// commands.
type RescanCmd struct {
	id         interface{}
	BeginBlock int32
	Addresses  map[string]struct{}
	EndBlock   int64 // TODO: switch this and btcdb.AllShas to int32
}

// Enforce that RescanCmd satisifies the btcjson.Cmd interface.
var _ btcjson.Cmd = &RescanCmd{}

// NewRescanCmd creates a new RescanCmd, parsing the optional
// arguments optArgs which may either be empty or a single upper
// block height.
func NewRescanCmd(id interface{}, begin int32, addresses map[string]struct{},
	optArgs ...int64) (*RescanCmd, error) {

	// Optional parameters set to their defaults.
	end := btcdb.AllShas

	if len(optArgs) > 0 {
		if len(optArgs) > 1 {
			return nil, btcjson.ErrTooManyOptArgs
		}
		end = optArgs[0]
	}

	return &RescanCmd{
		id:         id,
		BeginBlock: begin,
		Addresses:  addresses,
		EndBlock:   end,
	}, nil
}

// parseRescanCmd parses a RawCmd into a concrete type satisifying
// the btcjson.Cmd interface.  This is used when registering the custom
// command with the btcjson parser.
func parseRescanCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) < 2 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	begin, ok := r.Params[0].(float64)
	if !ok {
		return nil, errors.New("first parameter must be a number")
	}
	iaddrs, ok := r.Params[1].(map[string]interface{})
	if !ok {
		return nil, errors.New("second parameter must be a JSON object")
	}
	addresses := make(map[string]struct{}, len(iaddrs))
	for addr := range iaddrs {
		addresses[addr] = struct{}{}
	}
	params := make([]int64, len(r.Params[2:]))
	for i, val := range r.Params[2:] {
		fval, ok := val.(float64)
		if !ok {
			return nil, errors.New("optional parameters must " +
				"be be numbers")
		}
		params[i] = int64(fval)
	}

	return NewRescanCmd(r.Id, int32(begin), addresses, params...)
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *RescanCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *RescanCmd) Method() string {
	return "rescan"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *RescanCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "rescan",
		Id:      cmd.id,
		Params: []interface{}{
			cmd.BeginBlock,
			cmd.Addresses,
		},
	}

	if cmd.EndBlock != btcdb.AllShas {
		raw.Params = append(raw.Params, cmd.EndBlock)
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *RescanCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseRescanCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*RescanCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// NotifyNewTXsCmd is a type handling custom marshaling and
// unmarshaling of notifynewtxs JSON websocket extension
// commands.
type NotifyNewTXsCmd struct {
	id        interface{}
	Addresses []string
}

// Enforce that NotifyNewTXsCmd satisifies the btcjson.Cmd interface.
var _ btcjson.Cmd = &NotifyNewTXsCmd{}

// NewNotifyNewTXsCmd creates a new RescanCmd.
func NewNotifyNewTXsCmd(id interface{}, addresses []string) *NotifyNewTXsCmd {
	return &NotifyNewTXsCmd{
		id:        id,
		Addresses: addresses,
	}
}

// parseNotifyNewTXsCmd parses a NotifyNewTXsCmd into a concrete type
// satisifying the btcjson.Cmd interface.  This is used when registering
// the custom command with the btcjson parser.
func parseNotifyNewTXsCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 1 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	iaddrs, ok := r.Params[0].([]interface{})
	if !ok {
		return nil, errors.New("first parameter must be a JSON array")
	}
	addresses := make([]string, len(iaddrs))
	for i := range iaddrs {
		addr, ok := iaddrs[i].(string)
		if !ok {
			return nil, errors.New("first parameter must be an " +
				"array of strings")
		}
		addresses[i] = addr
	}

	return NewNotifyNewTXsCmd(r.Id, addresses), nil
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *NotifyNewTXsCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *NotifyNewTXsCmd) Method() string {
	return "notifynewtxs"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *NotifyNewTXsCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "notifynewtxs",
		Id:      cmd.id,
		Params: []interface{}{
			cmd.Addresses,
		},
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *NotifyNewTXsCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseNotifyNewTXsCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*NotifyNewTXsCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// NotifySpentCmd is a type handling custom marshaling and
// unmarshaling of notifyspent JSON websocket extension
// commands.
type NotifySpentCmd struct {
	id interface{}
	*btcwire.OutPoint
}

// Enforce that NotifySpentCmd satisifies the btcjson.Cmd interface.
var _ btcjson.Cmd = &NotifySpentCmd{}

// NewNotifySpentCmd creates a new RescanCmd.
func NewNotifySpentCmd(id interface{}, op *btcwire.OutPoint) *NotifySpentCmd {
	return &NotifySpentCmd{
		id:       id,
		OutPoint: op,
	}
}

// parseNotifySpentCmd parses a NotifySpentCmd into a concrete type
// satisifying the btcjson.Cmd interface.  This is used when registering
// the custom command with the btcjson parser.
func parseNotifySpentCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 2 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	hashStr, ok := r.Params[0].(string)
	if !ok {
		return nil, errors.New("first parameter must be a string")
	}
	hash, err := btcwire.NewShaHashFromStr(hashStr)
	if err != nil {
		return nil, errors.New("first parameter is not a valid " +
			"hash string")
	}
	idx, ok := r.Params[1].(float64)
	if !ok {
		return nil, errors.New("second parameter is not a number")
	}
	if idx < 0 {
		return nil, errors.New("second parameter cannot be negative")
	}

	cmd := NewNotifySpentCmd(r.Id, btcwire.NewOutPoint(hash, uint32(idx)))
	return cmd, nil
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *NotifySpentCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *NotifySpentCmd) Method() string {
	return "notifyspent"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *NotifySpentCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "notifyspent",
		Id:      cmd.id,
		Params: []interface{}{
			cmd.OutPoint.Hash.String(),
			cmd.OutPoint.Index,
		},
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *NotifySpentCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseNotifySpentCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*NotifySpentCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// CreateEncryptedWalletCmd is a type handling custom
// marshaling and unmarshaling of createencryptedwallet
// JSON websocket extension commands.
type CreateEncryptedWalletCmd struct {
	id          interface{}
	Account     string
	Description string
	Passphrase  string
}

// Enforce that CreateEncryptedWalletCmd satisifies the btcjson.Cmd
// interface.
var _ btcjson.Cmd = &CreateEncryptedWalletCmd{}

// NewCreateEncryptedWalletCmd creates a new CreateEncryptedWalletCmd.
func NewCreateEncryptedWalletCmd(id interface{},
	account, description, passphrase string) *CreateEncryptedWalletCmd {

	return &CreateEncryptedWalletCmd{
		id:          id,
		Account:     account,
		Description: description,
		Passphrase:  passphrase,
	}
}

// parseCreateEncryptedWalletCmd parses a CreateEncryptedWalletCmd
// into a concrete type satisifying the btcjson.Cmd interface.
// This is used when registering the custom command with the btcjson
// parser.
func parseCreateEncryptedWalletCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 3 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	account, ok := r.Params[0].(string)
	if !ok {
		return nil, errors.New("first parameter must be a string")
	}
	description, ok := r.Params[1].(string)
	if !ok {
		return nil, errors.New("second parameter is not a string")
	}
	passphrase, ok := r.Params[2].(string)
	if !ok {
		return nil, errors.New("third parameter is not a string")
	}

	cmd := NewCreateEncryptedWalletCmd(r.Id, account, description,
		passphrase)
	return cmd, nil
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *CreateEncryptedWalletCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *CreateEncryptedWalletCmd) Method() string {
	return "createencryptedwallet"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *CreateEncryptedWalletCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "createencryptedwallet",
		Id:      cmd.id,
		Params: []interface{}{
			cmd.Account,
			cmd.Description,
			cmd.Passphrase,
		},
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *CreateEncryptedWalletCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseCreateEncryptedWalletCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*CreateEncryptedWalletCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// GetBalancesCmd is a type handling custom marshaling and
// unmarshaling of getbalances JSON websocket extension commands.
type GetBalancesCmd struct {
	id interface{}
}

// Enforce that GetBalancesCmd satisifies the btcjson.Cmd
// interface.
var _ btcjson.Cmd = &GetBalancesCmd{}

// NewGetBalancesCmd creates a new GetBalancesCmd.
func NewGetBalancesCmd(id interface{}) *GetBalancesCmd {
	return &GetBalancesCmd{id: id}
}

// parseGetBalancesCmd parses a GetBalancesCmd into a concrete
// type satisifying the btcjson.Cmd interface.  This is used when
// registering the custom command with the btcjson parser.
func parseGetBalancesCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 0 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	return NewGetBalancesCmd(r.Id), nil
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *GetBalancesCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *GetBalancesCmd) Method() string {
	return "getbalances"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *GetBalancesCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "getbalances",
		Id:      cmd.id,
		Params:  []interface{}{},
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *GetBalancesCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseGetBalancesCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*GetBalancesCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// WalletIsLockedCmd is a type handling custom marshaling and
// unmarshaling of walletislocked JSON websocket extension commands.
type WalletIsLockedCmd struct {
	id      interface{}
	Account string
}

// Enforce that WalletIsLockedCmd satisifies the btcjson.Cmd
// interface.
var _ btcjson.Cmd = &WalletIsLockedCmd{}

// NewWalletIsLockedCmd creates a new WalletIsLockedCmd.
func NewWalletIsLockedCmd(id interface{},
	optArgs ...string) (*WalletIsLockedCmd, error) {

	// Optional arguments set to their default values.
	account := ""

	if len(optArgs) > 1 {
		return nil, btcjson.ErrInvalidParams
	}

	if len(optArgs) == 1 {
		account = optArgs[0]
	}

	return &WalletIsLockedCmd{
		id:      id,
		Account: account,
	}, nil
}

// parseWalletIsLockedCmd parses a WalletIsLockedCmd into a concrete
// type satisifying the btcjson.Cmd interface.  This is used when
// registering the custom command with the btcjson parser.
func parseWalletIsLockedCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) > 1 {
		return nil, btcjson.ErrInvalidParams
	}

	if len(r.Params) == 0 {
		return NewWalletIsLockedCmd(r.Id)
	}

	account, ok := r.Params[0].(string)
	if !ok {
		return nil, errors.New("account must be a string")
	}
	return NewWalletIsLockedCmd(r.Id, account)
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *WalletIsLockedCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *WalletIsLockedCmd) Method() string {
	return "walletislocked"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *WalletIsLockedCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "walletislocked",
		Id:      cmd.id,
		Params:  []interface{}{},
	}

	if cmd.Account != "" {
		raw.Params = append(raw.Params, cmd.Account)
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *WalletIsLockedCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseWalletIsLockedCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*WalletIsLockedCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}

// ListAllTransactionsCmd is a type handling custom marshaling and
// unmarshaling of listalltransactions JSON websocket extension commands.
type ListAllTransactionsCmd struct {
	id      interface{}
	Account string
}

// Enforce that ListAllTransactionsCmd satisifies the btcjson.Cmd
// interface.
var _ btcjson.Cmd = &ListAllTransactionsCmd{}

// NewListAllTransactionsCmd creates a new ListAllTransactionsCmd.
func NewListAllTransactionsCmd(id interface{},
	optArgs ...string) (*ListAllTransactionsCmd, error) {

	// Optional arguments set to their default values.
	account := ""

	if len(optArgs) > 1 {
		return nil, btcjson.ErrInvalidParams
	}

	if len(optArgs) == 1 {
		account = optArgs[0]
	}

	return &ListAllTransactionsCmd{
		id:      id,
		Account: account,
	}, nil
}

// parseListAllTransactionsCmd parses a ListAllTransactionsCmd into a concrete
// type satisifying the btcjson.Cmd interface.  This is used when
// registering the custom command with the btcjson parser.
func parseListAllTransactionsCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) > 1 {
		return nil, btcjson.ErrInvalidParams
	}

	if len(r.Params) == 0 {
		return NewListAllTransactionsCmd(r.Id)
	}

	account, ok := r.Params[0].(string)
	if !ok {
		return nil, errors.New("account must be a string")
	}
	return NewListAllTransactionsCmd(r.Id, account)
}

// Id satisifies the Cmd interface by returning the ID of the command.
func (cmd *ListAllTransactionsCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the RPC method.
func (cmd *ListAllTransactionsCmd) Method() string {
	return "listalltransactions"
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *ListAllTransactionsCmd) MarshalJSON() ([]byte, error) {
	// Fill a RawCmd and marshal.
	raw := btcjson.RawCmd{
		Jsonrpc: "1.0",
		Method:  "listalltransactions",
		Id:      cmd.id,
		Params:  []interface{}{},
	}

	if cmd.Account != "" {
		raw.Params = append(raw.Params, cmd.Account)
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *ListAllTransactionsCmd) UnmarshalJSON(b []byte) error {
	// Unmarshal into a RawCmd.
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	newCmd, err := parseListAllTransactionsCmd(&r)
	if err != nil {
		return err
	}

	concreteCmd, ok := newCmd.(*ListAllTransactionsCmd)
	if !ok {
		return btcjson.ErrInternal
	}
	*cmd = *concreteCmd
	return nil
}
