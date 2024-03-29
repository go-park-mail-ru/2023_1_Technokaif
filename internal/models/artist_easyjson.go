// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjson75485a89DecodeGithubComGoParkMailRu20231TechnokaifInternalModels(in *jlexer.Lexer, out *ArtistTransfers) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ArtistTransfers, 0, 1)
			} else {
				*out = ArtistTransfers{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 ArtistTransfer
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson75485a89EncodeGithubComGoParkMailRu20231TechnokaifInternalModels(out *jwriter.Writer, in ArtistTransfers) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ArtistTransfers) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson75485a89EncodeGithubComGoParkMailRu20231TechnokaifInternalModels(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ArtistTransfers) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson75485a89DecodeGithubComGoParkMailRu20231TechnokaifInternalModels(l, v)
}
func easyjson75485a89DecodeGithubComGoParkMailRu20231TechnokaifInternalModels1(in *jlexer.Lexer, out *ArtistTransfer) {
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
		case "name":
			out.Name = string(in.String())
		case "isLiked":
			out.IsLiked = bool(in.Bool())
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
func easyjson75485a89EncodeGithubComGoParkMailRu20231TechnokaifInternalModels1(out *jwriter.Writer, in ArtistTransfer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"isLiked\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsLiked))
	}
	{
		const prefix string = ",\"cover\":"
		out.RawString(prefix)
		out.String(string(in.AvatarSrc))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ArtistTransfer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson75485a89EncodeGithubComGoParkMailRu20231TechnokaifInternalModels1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ArtistTransfer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson75485a89DecodeGithubComGoParkMailRu20231TechnokaifInternalModels1(l, v)
}
