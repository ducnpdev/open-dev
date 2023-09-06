package main

// func TestNewBlake2b256(t *testing.T) {
// 	for i, tt := range []struct {
// 		in  []byte
// 		out string
// 	}{
// 		{[]byte(""), "0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"},
// 		{[]byte("abc"), "bddd813c634239723171ef3fee98579b94964e3bb1cb3e427262c8c068d52319"},
// 		{[]byte("hello"), "324dcf027dd4a30a932c441f365a25e86b173defa4b8e58948253471b81b72cf"},
// 	} {
// 		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
// 			result := NewBlake2b256(tt.in)
// 			if hex.EncodeToString(result) != tt.out {
// 				t.Errorf("want %v; got %v", tt.out, hex.EncodeToString(result))
// 			}
// 		})
// 	}
// }

// func TestNewBlake2b512(t *testing.T) {
// 	for i, tt := range []struct {
// 		in  []byte
// 		out string
// 	}{
// 		{[]byte(""), "786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce"},
// 		{[]byte("abc"), "ba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d17d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923"},
// 		{[]byte("hello"), "e4cfa39a3d37be31c59609e807970799caa68a19bfaa15135f165085e01d41a65ba1e1b146aeb6bd0092b49eac214c103ccfa3a365954bbbe52f74a2b3620c94"},
// 	} {
// 		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
// 			result := NewBlake2b512(tt.in)
// 			if hex.EncodeToString(result) != tt.out {
// 				t.Errorf("want %v; got %v", tt.out, hex.EncodeToString(result))
// 			}
// 		})
// 	}
// }
