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

func easyjson3b1bf41aDecodeGithubComGoParkMailRu20231TechnokaifInternalModels(in *jlexer.Lexer, out *PlaylistTransfers) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PlaylistTransfers, 0, 0)
			} else {
				*out = PlaylistTransfers{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 PlaylistTransfer
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
func easyjson3b1bf41aEncodeGithubComGoParkMailRu20231TechnokaifInternalModels(out *jwriter.Writer, in PlaylistTransfers) {
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
func (v PlaylistTransfers) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3b1bf41aEncodeGithubComGoParkMailRu20231TechnokaifInternalModels(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlaylistTransfers) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3b1bf41aDecodeGithubComGoParkMailRu20231TechnokaifInternalModels(l, v)
}
func easyjson3b1bf41aDecodeGithubComGoParkMailRu20231TechnokaifInternalModels1(in *jlexer.Lexer, out *PlaylistTransfer) {
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
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make(UserTransfers, 0, 0)
					} else {
						out.Users = UserTransfers{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v4 UserTransfer
					easyjson3b1bf41aDecodeGithubComGoParkMailRu20231TechnokaifInternalModels2(in, &v4)
					out.Users = append(out.Users, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "description":
			if in.IsNull() {
				in.Skip()
				out.Description = nil
			} else {
				if out.Description == nil {
					out.Description = new(string)
				}
				*out.Description = string(in.String())
			}
		case "isLiked":
			out.IsLiked = bool(in.Bool())
		case "cover":
			out.CoverSrc = string(in.String())
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
func easyjson3b1bf41aEncodeGithubComGoParkMailRu20231TechnokaifInternalModels1(out *jwriter.Writer, in PlaylistTransfer) {
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
		const prefix string = ",\"users\":"
		out.RawString(prefix)
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Users {
				if v5 > 0 {
					out.RawByte(',')
				}
				easyjson3b1bf41aEncodeGithubComGoParkMailRu20231TechnokaifInternalModels2(out, v6)
			}
			out.RawByte(']')
		}
	}
	if in.Description != nil {
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(*in.Description))
	}
	{
		const prefix string = ",\"isLiked\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsLiked))
	}
	if in.CoverSrc != "" {
		const prefix string = ",\"cover\":"
		out.RawString(prefix)
		out.String(string(in.CoverSrc))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlaylistTransfer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3b1bf41aEncodeGithubComGoParkMailRu20231TechnokaifInternalModels1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlaylistTransfer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3b1bf41aDecodeGithubComGoParkMailRu20231TechnokaifInternalModels1(l, v)
}
func easyjson3b1bf41aDecodeGithubComGoParkMailRu20231TechnokaifInternalModels2(in *jlexer.Lexer, out *UserTransfer) {
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
		case "username":
			out.Username = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "firstName":
			out.FirstName = string(in.String())
		case "lastName":
			out.LastName = string(in.String())
		case "birthDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.BirthDate).UnmarshalJSON(data))
			}
		case "avatarSrc":
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
func easyjson3b1bf41aEncodeGithubComGoParkMailRu20231TechnokaifInternalModels2(out *jwriter.Writer, in UserTransfer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"firstName\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"lastName\":"
		out.RawString(prefix)
		out.String(string(in.LastName))
	}
	if true {
		const prefix string = ",\"birthDate\":"
		out.RawString(prefix)
		out.Raw((in.BirthDate).MarshalJSON())
	}
	if in.AvatarSrc != "" {
		const prefix string = ",\"avatarSrc\":"
		out.RawString(prefix)
		out.String(string(in.AvatarSrc))
	}
	out.RawByte('}')
}
