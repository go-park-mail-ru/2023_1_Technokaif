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

func easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp(in *jlexer.Lexer, out *albumLikeResponse) {
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
func easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp(out *jwriter.Writer, in albumLikeResponse) {
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
func (v albumLikeResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *albumLikeResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp(l, v)
}
func easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp1(in *jlexer.Lexer, out *albumDeleteResponse) {
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
func easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp1(out *jwriter.Writer, in albumDeleteResponse) {
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
func (v albumDeleteResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *albumDeleteResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp1(l, v)
}
func easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp2(in *jlexer.Lexer, out *albumCreateResponse) {
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
func easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp2(out *jwriter.Writer, in albumCreateResponse) {
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
func (v albumCreateResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *albumCreateResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp2(l, v)
}
func easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp3(in *jlexer.Lexer, out *albumCreateInput) {
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
		case "artists":
			if in.IsNull() {
				in.Skip()
				out.ArtistsID = nil
			} else {
				in.Delim('[')
				if out.ArtistsID == nil {
					if !in.IsDelim(']') {
						out.ArtistsID = make([]uint32, 0, 16)
					} else {
						out.ArtistsID = []uint32{}
					}
				} else {
					out.ArtistsID = (out.ArtistsID)[:0]
				}
				for !in.IsDelim(']') {
					var v1 uint32
					v1 = uint32(in.Uint32())
					out.ArtistsID = append(out.ArtistsID, v1)
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
func easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp3(out *jwriter.Writer, in albumCreateInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"artists\":"
		out.RawString(prefix)
		if in.ArtistsID == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.ArtistsID {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.Uint32(uint32(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		if in.Description == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Description))
		}
	}
	{
		const prefix string = ",\"cover\":"
		out.RawString(prefix)
		out.String(string(in.CoverSrc))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v albumCreateInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBc29487EncodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *albumCreateInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBc29487DecodeGithubComGoParkMailRu20231TechnokaifInternalPkgAlbumDeliveryHttp3(l, v)
}
