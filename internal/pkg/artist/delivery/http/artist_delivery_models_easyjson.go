// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package http

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp(in *jlexer.Lexer, out *artistLikeResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "status":
			out.Status = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp(out *jwriter.Writer, in artistLikeResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v artistLikeResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *artistLikeResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp(l, v)
}
func easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp1(in *jlexer.Lexer, out *artistDeleteResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "status":
			out.Status = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp1(out *jwriter.Writer, in artistDeleteResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v artistDeleteResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *artistDeleteResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp1(l, v)
}
func easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp2(in *jlexer.Lexer, out *artistCreateResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = uint32(in.Uint32())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp2(out *jwriter.Writer, in artistCreateResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v artistCreateResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *artistCreateResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp2(l, v)
}
func easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp3(in *jlexer.Lexer, out *artistCreateInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "cover":
			out.AvatarSrc = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp3(out *jwriter.Writer, in artistCreateInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"cover\":"
		out.RawString(prefix)
		out.String(string(in.AvatarSrc))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v artistCreateInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11ae29f9EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *artistCreateInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11ae29f9DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgArtistDeliveryHttp3(l, v)
}
