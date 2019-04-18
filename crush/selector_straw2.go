package crush

//"log"

type Straw2Selector struct {
	ItemWeights map[Node]float32
}

func NewStraw2Selector(n Node) *Straw2Selector {
	s := new(Straw2Selector)
	s.ItemWeights = make(map[Node]float32)
	if n.IsLeaf() {
		return s
	}

	for _, child := range n.GetChildren() {
		s.ItemWeights[child] = child.GetWeight()
	}

	return s
}

func (s *Straw2Selector) Select(input uint32, round uint32) Node {
	var result Node
	var i uint32
	var highDraw int64

	for child, weight := range s.ItemWeights {
		draw := lnScore(child, weight, input, round)
		if i == 0 || draw > highDraw {
			highDraw = draw
			result = child
		}
		i++
	}

	if result == nil {
		panic("Illegal state")
	}
	return result
}

//compute 2^44*log2(input+1)
func crushLn(xin uint32) uint64 {
	var x, x1 uint32
	var iexpon, index1, index2 int
	var RH, LH, LL, xl64, result uint64

	x = xin + 1

	/* normalize input */
	iexpon = 15
	for 0 == (x & 0x18000) {
		x <<= 1
		iexpon--
	}

	index1 = int((x >> 8) << 1)
	/* RH ~ 2^56/index1 */
	RH = __RH_LH_tbl[index1-256]
	/* LH ~ 2^48 * log2(index1/256) */
	LH = __RH_LH_tbl[index1+1-256]

	/* RH*x ~ 2^48 * (2^15 + xf), xf<2^8 */
	xl64 = uint64(x) * RH
	xl64 >>= 48
	x1 = uint32(xl64)

	result = uint64(iexpon)
	result <<= (12 + 32)

	index2 = int(x1 & 0xff)
	/* LL ~ 2^48*log2(1.0+index2/2^15) */
	LL = uint64(__LL_tbl[index2])

	LH = LH + LL

	LH >>= (48 - 12 - 32)
	result += LH

	return result
}

/*
 * RH_LH_tbl[2*k] = 2^48/(1.0+k/128.0)
 * RH_LH_tbl[2*k+1] = 2^48*log2(1.0+k/128.0)
 */
var __RH_LH_tbl []uint64 = []uint64{
	0x0001000000000000, 0x0000000000000000, 0x0000fe03f80fe040, 0x000002dfca16dde1,
	0x0000fc0fc0fc0fc1, 0x000005b9e5a170b4, 0x0000fa232cf25214, 0x0000088e68ea899a,
	0x0000f83e0f83e0f9, 0x00000b5d69bac77e, 0x0000f6603d980f67, 0x00000e26fd5c8555,
	0x0000f4898d5f85bc, 0x000010eb389fa29f, 0x0000f2b9d6480f2c, 0x000013aa2fdd27f1,
	0x0000f0f0f0f0f0f1, 0x00001663f6fac913, 0x0000ef2eb71fc435, 0x00001918a16e4633,
	0x0000ed7303b5cc0f, 0x00001bc84240adab, 0x0000ebbdb2a5c162, 0x00001e72ec117fa5,
	0x0000ea0ea0ea0ea1, 0x00002118b119b4f3, 0x0000e865ac7b7604, 0x000023b9a32eaa56,
	0x0000e6c2b4481cd9, 0x00002655d3c4f15c, 0x0000e525982af70d, 0x000028ed53f307ee,
	0x0000e38e38e38e39, 0x00002b803473f7ad, 0x0000e1fc780e1fc8, 0x00002e0e85a9de04,
	0x0000e070381c0e08, 0x0000309857a05e07, 0x0000dee95c4ca038, 0x0000331dba0efce1,
	0x0000dd67c8a60dd7, 0x0000359ebc5b69d9, 0x0000dbeb61eed19d, 0x0000381b6d9bb29b,
	0x0000da740da740db, 0x00003a93dc9864b2, 0x0000d901b2036407, 0x00003d0817ce9cd4,
	0x0000d79435e50d7a, 0x00003f782d7204d0, 0x0000d62b80d62b81, 0x000041e42b6ec0c0,
	0x0000d4c77b03531e, 0x0000444c1f6b4c2d, 0x0000d3680d3680d4, 0x000046b016ca47c1,
	0x0000d20d20d20d21, 0x000049101eac381c, 0x0000d0b69fcbd259, 0x00004b6c43f1366a,
	0x0000cf6474a8819f, 0x00004dc4933a9337, 0x0000ce168a772509, 0x0000501918ec6c11,
	0x0000cccccccccccd, 0x00005269e12f346e, 0x0000cb8727c065c4, 0x000054b6f7f1325a,
	0x0000ca4587e6b750, 0x0000570068e7ef5a, 0x0000c907da4e8712, 0x000059463f919dee,
	0x0000c7ce0c7ce0c8, 0x00005b8887367433, 0x0000c6980c6980c7, 0x00005dc74ae9fbec,
	0x0000c565c87b5f9e, 0x00006002958c5871, 0x0000c4372f855d83, 0x0000623a71cb82c8,
	0x0000c30c30c30c31, 0x0000646eea247c5c, 0x0000c1e4bbd595f7, 0x000066a008e4788c,
	0x0000c0c0c0c0c0c1, 0x000068cdd829fd81, 0x0000bfa02fe80bfb, 0x00006af861e5fc7d,
	0x0000be82fa0be830, 0x00006d1fafdce20a, 0x0000bd6910470767, 0x00006f43cba79e40,
	0x0000bc52640bc527, 0x00007164beb4a56d, 0x0000bb3ee721a54e, 0x000073829248e961,
	0x0000ba2e8ba2e8bb, 0x0000759d4f80cba8, 0x0000b92143fa36f6, 0x000077b4ff5108d9,
	0x0000b81702e05c0c, 0x000079c9aa879d53, 0x0000b70fbb5a19bf, 0x00007bdb59cca388,
	0x0000b60b60b60b61, 0x00007dea15a32c1b, 0x0000b509e68a9b95, 0x00007ff5e66a0ffe,
	0x0000b40b40b40b41, 0x000081fed45cbccb, 0x0000b30f63528918, 0x00008404e793fb81,
	0x0000b21642c8590c, 0x000086082806b1d5, 0x0000b11fd3b80b12, 0x000088089d8a9e47,
	0x0000b02c0b02c0b1, 0x00008a064fd50f2a, 0x0000af3addc680b0, 0x00008c01467b94bb,
	0x0000ae4c415c9883, 0x00008df988f4ae80, 0x0000ad602b580ad7, 0x00008fef1e987409,
	0x0000ac7691840ac8, 0x000091e20ea1393e, 0x0000ab8f69e2835a, 0x000093d2602c2e5f,
	0x0000aaaaaaaaaaab, 0x000095c01a39fbd6, 0x0000a9c84a47a080, 0x000097ab43af59f9,
	0x0000a8e83f5717c1, 0x00009993e355a4e5, 0x0000a80a80a80a81, 0x00009b79ffdb6c8b,
	0x0000a72f0539782a, 0x00009d5d9fd5010b, 0x0000a655c4392d7c, 0x00009f3ec9bcfb80,
	0x0000a57eb50295fb, 0x0000a11d83f4c355, 0x0000a4a9cf1d9684, 0x0000a2f9d4c51039,
	0x0000a3d70a3d70a4, 0x0000a4d3c25e68dc, 0x0000a3065e3fae7d, 0x0000a6ab52d99e76,
	0x0000a237c32b16d0, 0x0000a8808c384547, 0x0000a16b312ea8fd, 0x0000aa5374652a1c,
	0x0000a0a0a0a0a0a1, 0x0000ac241134c4e9, 0x00009fd809fd80a0, 0x0000adf26865a8a1,
	0x00009f1165e72549, 0x0000afbe7fa0f04d, 0x00009e4cad23dd60, 0x0000b1885c7aa982,
	0x00009d89d89d89d9, 0x0000b35004723c46, 0x00009cc8e160c3fc, 0x0000b5157cf2d078,
	0x00009c09c09c09c1, 0x0000b6d8cb53b0ca, 0x00009b4c6f9ef03b, 0x0000b899f4d8ab63,
	0x00009a90e7d95bc7, 0x0000ba58feb2703a, 0x000099d722dabde6, 0x0000bc15edfeed32,
	0x0000991f1a515886, 0x0000bdd0c7c9a817, 0x00009868c809868d, 0x0000bf89910c1678,
	0x000097b425ed097c, 0x0000c1404eadf383, 0x000097012e025c05, 0x0000c2f5058593d9,
	0x0000964fda6c0965, 0x0000c4a7ba58377c, 0x000095a02568095b, 0x0000c65871da59dd,
	0x000094f2094f2095, 0x0000c80730b00016, 0x0000944580944581, 0x0000c9b3fb6d0559,
	0x0000939a85c4093a, 0x0000cb5ed69565af, 0x000092f113840498, 0x0000cd07c69d8702,
	0x0000924924924925, 0x0000ceaecfea8085, 0x000091a2b3c4d5e7, 0x0000d053f6d26089,
	0x000090fdbc090fdc, 0x0000d1f73f9c70c0, 0x0000905a38633e07, 0x0000d398ae817906,
	0x00008fb823ee08fc, 0x0000d53847ac00a6, 0x00008f1779d9fdc4, 0x0000d6d60f388e41,
	0x00008e78356d1409, 0x0000d8720935e643, 0x00008dda5202376a, 0x0000da0c39a54804,
	0x00008d3dcb08d3dd, 0x0000dba4a47aa996, 0x00008ca29c046515, 0x0000dd3b4d9cf24b,
	0x00008c08c08c08c1, 0x0000ded038e633f3, 0x00008b70344a139c, 0x0000e0636a23e2ee,
	0x00008ad8f2fba939, 0x0000e1f4e5170d02, 0x00008a42f870566a, 0x0000e384ad748f0e,
	0x000089ae4089ae41, 0x0000e512c6e54998, 0x0000891ac73ae982, 0x0000e69f35065448,
	0x0000888888888889, 0x0000e829fb693044, 0x000087f78087f781, 0x0000e9b31d93f98e,
	0x00008767ab5f34e5, 0x0000eb3a9f019750, 0x000086d905447a35, 0x0000ecc08321eb30,
	0x0000864b8a7de6d2, 0x0000ee44cd59ffab, 0x000085bf37612cef, 0x0000efc781043579,
	0x0000853408534086, 0x0000f148a170700a, 0x000084a9f9c8084b, 0x0000f2c831e44116,
	0x0000842108421085, 0x0000f446359b1353, 0x0000839930523fbf, 0x0000f5c2afc65447,
	0x000083126e978d50, 0x0000f73da38d9d4a, 0x0000828cbfbeb9a1, 0x0000f8b7140edbb1,
	0x0000820820820821, 0x0000fa2f045e7832, 0x000081848da8faf1, 0x0000fba577877d7d,
	0x0000810204081021, 0x0000fd1a708bbe11, 0x0000808080808081, 0x0000fe8df263f957,
	0x0000800000000000, 0x0000ffff00000000,
}

/*
 * LL_tbl[k] = 2^48*log2(1.0+k/2^15)
 */
var __LL_tbl []int64 = []int64{
	0x0000000000000000, 0x00000002e2a60a00, 0x000000070cb64ec5, 0x00000009ef50ce67,
	0x0000000cd1e588fd, 0x0000000fb4747e9c, 0x0000001296fdaf5e, 0x0000001579811b58,
	0x000000185bfec2a1, 0x0000001b3e76a552, 0x0000001e20e8c380, 0x0000002103551d43,
	0x00000023e5bbb2b2, 0x00000026c81c83e4, 0x00000029aa7790f0, 0x0000002c8cccd9ed,
	0x0000002f6f1c5ef2, 0x0000003251662017, 0x0000003533aa1d71, 0x0000003815e8571a,
	0x0000003af820cd26, 0x0000003dda537fae, 0x00000040bc806ec8, 0x000000439ea79a8c,
	0x0000004680c90310, 0x0000004962e4a86c, 0x0000004c44fa8ab6, 0x0000004f270aaa06,
	0x0000005209150672, 0x00000054eb19a013, 0x00000057cd1876fd, 0x0000005aaf118b4a,
	0x0000005d9104dd0f, 0x0000006072f26c64, 0x0000006354da3960, 0x0000006636bc441a,
	0x0000006918988ca8, 0x0000006bfa6f1322, 0x0000006edc3fd79f, 0x00000071be0ada35,
	0x000000749fd01afd, 0x00000077818f9a0c, 0x0000007a6349577a, 0x0000007d44fd535e,
	0x0000008026ab8dce, 0x00000083085406e3, 0x00000085e9f6beb2, 0x00000088cb93b552,
	0x0000008bad2aeadc, 0x0000008e8ebc5f65, 0x0000009170481305, 0x0000009451ce05d3,
	0x00000097334e37e5, 0x0000009a14c8a953, 0x0000009cf63d5a33, 0x0000009fd7ac4a9d,
	0x000000a2b07f3458, 0x000000a59a78ea6a, 0x000000a87bd699fb, 0x000000ab5d2e8970,
	0x000000ae3e80b8e3, 0x000000b11fcd2869, 0x000000b40113d818, 0x000000b6e254c80a,
	0x000000b9c38ff853, 0x000000bca4c5690c, 0x000000bf85f51a4a, 0x000000c2671f0c26,
	0x000000c548433eb6, 0x000000c82961b211, 0x000000cb0a7a664d, 0x000000cdeb8d5b82,
	0x000000d0cc9a91c8, 0x000000d3ada20933, 0x000000d68ea3c1dd, 0x000000d96f9fbbdb,
	0x000000dc5095f744, 0x000000df31867430, 0x000000e2127132b5, 0x000000e4f35632ea,
	0x000000e7d43574e6, 0x000000eab50ef8c1, 0x000000ed95e2be90, 0x000000f076b0c66c,
	0x000000f35779106a, 0x000000f6383b9ca2, 0x000000f918f86b2a, 0x000000fbf9af7c1a,
	0x000000feda60cf88, 0x00000101bb0c658c, 0x000001049bb23e3c, 0x000001077c5259af,
	0x0000010a5cecb7fc, 0x0000010d3d81593a, 0x000001101e103d7f, 0x00000112fe9964e4,
	0x00000115df1ccf7e, 0x00000118bf9a7d64, 0x0000011ba0126ead, 0x0000011e8084a371,
	0x0000012160f11bc6, 0x000001244157d7c3, 0x0000012721b8d77f, 0x0000012a02141b10,
	0x0000012ce269a28e, 0x0000012fc2b96e0f, 0x00000132a3037daa, 0x000001358347d177,
	0x000001386386698c, 0x0000013b43bf45ff, 0x0000013e23f266e9, 0x00000141041fcc5e,
	0x00000143e4477678, 0x00000146c469654b, 0x00000149a48598f0, 0x0000014c849c117c,
	0x0000014f64accf08, 0x0000015244b7d1a9, 0x0000015524bd1976, 0x0000015804bca687,
	0x0000015ae4b678f2, 0x0000015dc4aa90ce, 0x00000160a498ee31, 0x0000016384819134,
	0x00000166646479ec, 0x000001694441a870, 0x0000016c24191cd7, 0x0000016df6ca19bd,
	0x00000171e3b6d7aa, 0x00000174c37d1e44, 0x00000177a33dab1c, 0x0000017a82f87e49,
	0x0000017d62ad97e2, 0x00000180425cf7fe, 0x00000182b07f3458, 0x0000018601aa8c19,
	0x00000188e148c046, 0x0000018bc0e13b52, 0x0000018ea073fd52, 0x000001918001065d,
	0x000001945f88568b, 0x000001973f09edf2, 0x0000019a1e85ccaa, 0x0000019cfdfbf2c8,
	0x0000019fdd6c6063, 0x000001a2bcd71593, 0x000001a59c3c126e, 0x000001a87b9b570b,
	0x000001ab5af4e380, 0x000001ae3a48b7e5, 0x000001b11996d450, 0x000001b3f8df38d9,
	0x000001b6d821e595, 0x000001b9b75eda9b, 0x000001bc96961803, 0x000001bf75c79de3,
	0x000001c254f36c51, 0x000001c534198365, 0x000001c81339e336, 0x000001caf2548bd9,
	0x000001cdd1697d67, 0x000001d0b078b7f5, 0x000001d38f823b9a, 0x000001d66e86086d,
	0x000001d94d841e86, 0x000001dc2c7c7df9, 0x000001df0b6f26df, 0x000001e1ea5c194e,
	0x000001e4c943555d, 0x000001e7a824db23, 0x000001ea8700aab5, 0x000001ed65d6c42b,
	0x000001f044a7279d, 0x000001f32371d51f, 0x000001f60236ccca, 0x000001f8e0f60eb3,
	0x000001fbbfaf9af3, 0x000001fe9e63719e, 0x000002017d1192cc, 0x000002045bb9fe94,
	0x000002073a5cb50d, 0x00000209c06e6212, 0x0000020cf791026a, 0x0000020fd622997c,
	0x00000212b07f3458, 0x000002159334a8d8, 0x0000021871b52150, 0x0000021b502fe517,
	0x0000021d6a73a78f, 0x000002210d144eee, 0x00000223eb7df52c, 0x00000226c9e1e713,
	0x00000229a84024bb, 0x0000022c23679b4e, 0x0000022f64eb83a8, 0x000002324338a51b,
	0x00000235218012a9, 0x00000237ffc1cc69, 0x0000023a2c3b0ea4, 0x0000023d13ee805b,
	0x0000024035e9221f, 0x00000243788faf25, 0x0000024656b4e735, 0x00000247ed646bfe,
	0x0000024c12ee3d98, 0x0000024ef1025c1a, 0x00000251cf10c799, 0x0000025492644d65,
	0x000002578b1c85ee, 0x0000025a6919d8f0, 0x0000025d13ee805b, 0x0000026025036716,
	0x0000026296453882, 0x00000265e0d62b53, 0x00000268beb701f3, 0x0000026b9c92265e,
	0x0000026d32f798a9, 0x00000271583758eb, 0x000002743601673b, 0x0000027713c5c3b0,
	0x00000279f1846e5f, 0x0000027ccf3d6761, 0x0000027e6580aecb, 0x000002828a9e44b3,
	0x0000028568462932, 0x00000287bdbf5255, 0x0000028b2384de4a, 0x0000028d13ee805b,
	0x0000029035e9221f, 0x0000029296453882, 0x0000029699bdfb61, 0x0000029902a37aab,
	0x0000029c54b864c9, 0x0000029deabd1083, 0x000002a20f9c0bb5, 0x000002a4c7605d61,
	0x000002a7bdbf5255, 0x000002a96056dafc, 0x000002ac3daf14ef, 0x000002af1b019eca,
	0x000002b296453882, 0x000002b5d022d80f, 0x000002b8fa471cb3, 0x000002ba9012e713,
	0x000002bd6d4901cc, 0x000002c04a796cf6, 0x000002c327a428a6, 0x000002c61a5e8f4c,
	0x000002c8e1e891f6, 0x000002cbbf023fc2, 0x000002ce9c163e6e, 0x000002d179248e13,
	0x000002d4562d2ec6, 0x000002d73330209d, 0x000002da102d63b0, 0x000002dced24f814,
}
