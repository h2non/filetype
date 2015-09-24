package matchers

var TypeMidi = newType("mid", "audio/midi")
var TypeMp3 = newType("mp3", "audio/mpeg")
var TypeM4a = newType("m4a", "audio/m4a")
var TypeOgg = newType("ogg", "audio/ogg")
var TypeFlac = newType("flac", "audio/x-flac")
var TypeWav = newType("wav", "audio/x-wav")

var Audio = Map{
	TypeMidi: Midi,
	TypeMp3:  Mp3,
	TypeM4a:  M4a,
	TypeOgg:  Ogg,
	TypeFlac: Flac,
	TypeWav:  Wav,
}

func Midi(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x4D && buf[1] == 0x54 &&
		buf[2] == 0x68 && buf[3] == 0x64
}

func Mp3(buf []byte, length int) bool {
	return length > 2 &&
		(buf[0] == 0x49 && buf[1] == 0x44 && buf[2] == 0x33) ||
		(buf[0] == 0xFF && buf[1] == 0xfb)
}

func M4a(buf []byte, length int) bool {
	return length > 10 &&
		(buf[4] == 0x66 && buf[5] == 0x74 && buf[6] == 0x79 &&
			buf[7] == 0x70 && buf[8] == 0x4D && buf[9] == 0x34 && buf[10] == 0x41) ||
		(buf[0] == 0x4D && buf[1] == 0x34 && buf[2] == 0x41 && buf[3] == 0x20)
}

func Ogg(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x4F && buf[1] == 0x67 &&
		buf[2] == 0x67 && buf[3] == 0x53
}

func Flac(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x66 && buf[1] == 0x4C &&
		buf[2] == 0x61 && buf[3] == 0x43
}

func Wav(buf []byte, length int) bool {
	return length > 11 &&
		buf[0] == 0x52 && buf[1] == 0x49 && buf[2] == 0x46 && buf[3] == 0x46 &&
		buf[8] == 0x57 && buf[9] == 0x41 && buf[10] == 0x56 && buf[11] == 0x45
}
