package main

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
	"github.com/ssttevee/go-av"
	"github.com/ssttevee/go-av/avformat"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// if ((ret = open_input_file(argv[1])) < 0)
// goto end;
// if ((ret = open_output_file(argv[2])) < 0)
// goto end;
// if ((ret = init_filters()) < 0)
// goto end;
// if (!(packet = av_packet_alloc()))

var ifmt_ctx *av.InputFormatContext
var istreams []*av.Stream
var ofmt_ctx *av.OutputFormatContext
var ostreams []*av.Stream

func openInputFile(filename string) {
	var err error
	ifmt_ctx, err = av.OpenInputFile(filename)
	must(errors.WithStack(err))
	for _, stream := range ifmt_ctx.Streams() {
		dec, err := av.FindDecoderCodecByID(stream.Codecpar().CodecID)
		must(errors.WithStack(err))

		codec_ctx, err := av.NewDecoderContext(dec, stream.Codecpar())
		must(errors.WithStack(err))

		err = codec_ctx.Open()
		must(errors.WithStack(err))
		istreams = append(istreams, stream)
	}
	avformat.Dump(ifmt_ctx.GetFormatContext(), 1, filename, 0)
	fmt.Println(reflect.TypeOf(ifmt_ctx.GetFormatContext()))
	// avutil
}

func openOutputFile(filename string) {
	var err error
	ofmt_ctx, err = av.NewFileOutputContext("mp4", filename)
	must(errors.WithStack(err))

	for _, istream := range istreams {
		enc, err := av.FindEncoderCodecByID(istream.Codecpar().CodecID)
		must(errors.WithStack(err))

		// fmt.Printf("%++v\n", istream.CodecParameters().Inner())
		// os.Exit(0)

		codec_ctx, err := av.NewEncoderContext(enc, istream.Codecpar())
		must(errors.WithStack(err))

		codec_ctx.TimeBase = istream.TimeBase
		fmt.Println("i", istream.TimeBase)

		err = codec_ctx.Open()
		must(errors.WithStack(err))

		stream := ofmt_ctx.NewStream(enc)
		stream.SetCodecpar(istream.Codecpar())
		ostreams = append(ostreams, stream)

	}

	avformat.Dump(ofmt_ctx.GetFormatContext(), 1, filename, 1)
	pkt := av.NewPacket()
	fmt.Println(ofmt_ctx.WritePacket(pkt))
	// // ofmt_ctx
	// ret := avformat.WriteHeader(ofmt_ctx.GetFormatContext(), nil)
	// if ret != 0 {
	// 	must(fmt.Errorf("failed to write header: %d", ret))
	// }
}

func main() {
	openInputFile("1.mp4")
	openOutputFile("out.mp4")
	// streams := format.Streams()

	// vc, err := av.FindEncoderCodecByName("libx264")
	// must(errors.WithStack(err))

	// vcc, err := av.NewEncoderContext(vc, nil)
	// must(errors.WithStack(err))

	// vcc.SendFrame()
	// format.ReadPacket()

	// n := 0
	// for format.Pb.EofReached != 1 {
	// 	pkt, err := format.ReadPacket()
	// 	if pkt == nil {
	// 		fmt.Println("EOF")
	// 	}
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		break
	// 	}
	// 	fmt.Printf("%d #%d: %.2f\n", n, pkt.StreamIndex, float64(pkt.Pts)*streams[pkt.StreamIndex].TimeBase.Float64())

	// 	pkt

	// 	n++
	// 	// time.Sleep(time.Millisecond * 500)
	// }

	// invoke gc to test finalizer
	runtime.GC()
}
