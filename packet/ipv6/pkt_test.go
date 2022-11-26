/*
 * Network packet analysis framework.
 *
 * Copyright (c) 2014, Alessandro Ghedini
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
 * IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
 * THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
 * PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package ipv6_test

import "bytes"
import "net"
import "testing"

import "github.com/scs-solution/go.pkt2/packet"
import "github.com/scs-solution/go.pkt2/packet/ipv4"
import "github.com/scs-solution/go.pkt2/packet/ipv6"

var test_simple = []byte{
	0x63, 0x0d, 0x5b, 0x0a, 0x00, 0x08, 0x11, 0x40, 0xfe, 0x80, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x4e, 0x72, 0xb9, 0xff, 0xfe, 0x54, 0xe5, 0x3d,
	0x07, 0x9a, 0x19, 0xb9, 0x11, 0x15, 0xed, 0x67, 0x99, 0xf5, 0xf0, 0x7a,
	0x66, 0x87, 0x5b, 0x0f,
}

var ipsrc_str = "fe80::4e72:b9ff:fe54:e53d"
var ipdst_str = "79a:19b9:1115:ed67:99f5:f07a:6687:5b0f"

func MakeTestSimple() *ipv6.Packet {
	return &ipv6.Packet{
		Version:  6,
		Class:    48,
		Label:    875274,
		Length:   8,
		NextHdr:  ipv4.UDP,
		HopLimit: 64,
		SrcAddr:  net.ParseIP(ipsrc_str),
		DstAddr:  net.ParseIP(ipdst_str),
	}
}

func TestPack(t *testing.T) {
	var b packet.Buffer
	b.Init(make([]byte, len(test_simple)))

	p := MakeTestSimple()

	err := p.Pack(&b)
	if err != nil {
		t.Fatalf("Error packing: %s", err)
	}

	if !bytes.Equal(test_simple, b.Buffer()) {
		t.Fatalf("Raw packet mismatch: %x", b.Buffer())
	}
}

func BenchmarkPack(bn *testing.B) {
	var b packet.Buffer
	b.Init(make([]byte, len(test_simple)))

	p := MakeTestSimple()

	for n := 0; n < bn.N; n++ {
		p.Pack(&b)
	}
}

func TestUnpack(t *testing.T) {
	var p ipv6.Packet

	cmp := MakeTestSimple()

	var b packet.Buffer
	b.Init(test_simple)

	err := p.Unpack(&b)
	if err != nil {
		t.Fatalf("Error unpacking: %s", err)
	}

	if !p.Equals(cmp) {
		t.Fatalf("Packet mismatch:\n%s\n%s", &p, cmp)
	}
}

func BenchmarkUnpack(bn *testing.B) {
	var p ipv6.Packet
	var b packet.Buffer

	for n := 0; n < bn.N; n++ {
		b.Init(test_simple)
		p.Unpack(&b)
	}
}
