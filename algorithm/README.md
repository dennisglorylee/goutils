**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [algorithm](#algorithm)
  * [circ2buffer_test.go](#circ2buffer_testgo)
  * [crc16-kermit_test.go](#crc16-kermit_testgo)
  * [crc16_test.go](#crc16_testgo)
  * [descartes_test.go](#descartes_testgo)
  * [xor_io_test.go](#xor_io_testgo)

<!-- tocstop -->

# algorithm
## circ2buffer_test.go
### TestCreateC2Buffer
```go

MakeC2Buffer(BLOCK_SIZE)
```
### TestWriteBlock
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
```
### TestWritingUnderCapacityGivesEmptyEvicted
```go

b := MakeC2Buffer(2)
b.Write([]byte{1, 2})

if len(b.Evicted()) != 0 {
	t.Fatal("Evicted should have been empty:", b.Evicted())
}
```
### TestWritingMultipleBytesWhenBufferIsNotFull
```go

b := MakeC2Buffer(3)
b.Write([]byte{1, 2})
b.Write([]byte{3, 4})

ev := b.Evicted()

if len(ev) != 1 || ev[0] != 1 {
	t.Fatal("Evicted should have been [1,]:", ev)
}
```
### TestEvictedRegession1
```go

b := MakeC2Buffer(4)

b.Write([]byte{7, 6})
b.Write([]byte{5, 1, 2})
b.Write([]byte{3, 4})

ev := b.Evicted()
if len(ev) != 2 || ev[0] != 6 || ev[1] != 5 {
	t.Fatalf("Unexpected evicted [6,5]: %v", ev)
}
```
### TestGetBlock
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)

block := b.GetBlock()

if len(block) != BLOCK_SIZE {
	t.Fatal("Wrong block size returned")
}

for i, by := range block {
	if byte(i) != by {
		t.Errorf("byte %v does not match", i)
	}
}
```
### TestWriteTwoBlocksGet
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
b.Write(incrementBlock2)

if bytes.Compare(b.GetBlock(), incrementBlock2) != 0 {
	t.Errorf("Get block did not return the right value: %s", b.GetBlock())
}
```
### TestWriteSingleByteGetSingleByte
```go

b := MakeC2Buffer(BLOCK_SIZE)
singleByte := []byte{0}
b.Write(singleByte)

if bytes.Compare(b.GetBlock(), singleByte) != 0 {
	t.Errorf("Get block did not return the right value: %s", b.GetBlock())
}
```
### TestWriteTwoBlocksGetEvicted
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
b.Write(incrementBlock2)

if bytes.Compare(b.Evicted(), incrementBlock) != 0 {
	t.Errorf("Evicted did not return the right value: %s", b.Evicted())
}
```
### TestWriteSingleByteReturnsSingleEvictedByte
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock2)
singleByte := []byte{0}

b.Write(singleByte)
e := b.Evicted()

if len(e) != 1 {
	t.Fatalf("Evicted length is not correct: %s", e)
}

if e[0] != byte(10) {
	t.Errorf("Evicted content is not correct: %s", e)
}
```
### TestTruncatingAfterWriting
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)

evicted := b.Truncate(2)

if len(evicted) != 2 {
	t.Fatalf("Truncate did not return expected evicted length: %v", evicted)
}

if evicted[0] != 0 || evicted[1] != 1 {
	t.Errorf("Unexpected content in evicted: %v", evicted)
}
```
### TestWritingAfterTruncating
```go

// test that after we truncate some content, the next operations
// on the buffer give us the expected results
b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
b.Truncate(4)

b.Write([]byte{34, 46})

block := b.GetBlock()

if len(block) != BLOCK_SIZE-2 {
	t.Fatalf(
		"Unexpected block length after truncation: %v (%v)",
		block,
		len(block),
	)
}

if bytes.Compare(block, []byte{4, 5, 6, 7, 8, 9, 34, 46}) != 0 {
	t.Errorf(
		"Unexpected block content after truncation: %v (%v)",
		block,
		len(block))
}
```
## crc16-kermit_test.go
### TestKermit
```go

a := Kermit([]byte("abcdefg汉字"))
b := Kermit([]byte("abcdefg汉字"))
if a != b {
	t.Error(a, b)
}
```
## crc16_test.go
### TestCrc16
```go

a := Crc16([]byte("abcdefg汉字"))
b := Crc16([]byte("abcdefg汉字"))
if a != b {
	t.Error(Crc16([]byte("abcdefg汉字")))
}
```
### TestCrc16s
```go

a := Crc16s("abcdefg汉字")
b := Crc16([]byte("abcdefg汉字"))

if a != b {
	t.Error(Crc16([]byte("abcdefg汉字")))
}
```
## descartes_test.go
### TestDescartes
```go

result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})

descartMap := make(map[string]bool)
for _, item := range result {
	if ok, _ := testDescartesStrContains([]string{"A", "B"}, item[0]); !ok {
		t.FailNow()
	}

	if ok, _ := testDescartesStrContains([]string{"1", "2", "3"}, item[1]); !ok {
		t.FailNow()
	}

	if ok, _ := testDescartesStrContains([]string{"a", "b", "c", "d"}, item[2]); !ok {
		t.FailNow()
	}
	descartMap[fmt.Sprint(item)] = true
}

if len(descartMap) != 24 {
	t.FailNow()
}
```
## xor_io_test.go
### TestXorIO
```go

data := []byte("1234567890abcdefhijklmn")

w := &bytes.Buffer{}

xw := NewXORWriter(w, testXorKey)
_, err := io.Copy(xw, bytes.NewReader(data))
if err != nil {
	t.Fatal(err)
}

cipherBs := w.Bytes()
xr := NewXORReader(bytes.NewReader(cipherBs), testXorKey)
rdata, err := ioutil.ReadAll(xr)
if err != nil {
	t.Fatal(err)
}

if bytes.Compare(data, rdata) != 0 {
	t.FailNow()
}
```
### cipherXor
```go

writerIndex := uint64(0)

f, err := os.Open(testXorOrigFilePath)
if err != nil {
	t.Fatal(err)
}
defer f.Close()

cf, err := os.OpenFile(testXorCrpytoFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
if err != nil {
	t.Fatal(err)
}
defer cf.Close()

w := NewXORWriterWithOffset(cf, testXorKey, &writerIndex)

_, err = io.Copy(w, f)
if err != nil {
	t.Fatal(err)
}
```
### TestDeCipherXor
```go

readerIndex := uint64(0)

cipherXor(t)

func() {
	f, err := os.Open(testXorCrpytoFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := NewXORReaderWithOffset(f, testXorKey, &readerIndex)

	rf, err := os.OpenFile(testXorRecoverFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer rf.Close()

	_, err = io.Copy(rf, r)
	if err != nil {
		t.Fatal(err)
	}
}()

bs1, _ := ioutil.ReadFile(testXorOrigFilePath)
bs2, _ := ioutil.ReadFile(testXorRecoverFilePath)

if bytes.Compare(bs1, bs2) != 0 {
	t.FailNow()
}
```
### TestXORReaderAt
```go

cipherXor(t)
func() {

	f, err := os.Open(testXorCrpytoFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rf, err := os.OpenFile(testXorRecoverFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer rf.Close()

	fileInfo, err := os.Stat(testXorCrpytoFilePath)
	if err != nil {
		return
	}
	size := fileInfo.Size()

	for offset := int64(0); offset < size; {
		rdsize := int64(rand.Intn(int(size) / 2))
		if offset+rdsize > size {
			rdsize = size - offset
		}

		t.Logf("read section offset: %v, size: %v\n", offset, rdsize)
		sr := io.NewSectionReader(NewXORReaderAt(f, testXorKey), offset, rdsize)

		_, err = io.Copy(rf, sr)
		if err != nil {
			t.Fatal(err)
		}

		offset += rdsize
	}
}()

bs1, _ := ioutil.ReadFile(testXorOrigFilePath)
bs2, _ := ioutil.ReadFile(testXorRecoverFilePath)

if bytes.Compare(bs1, bs2) != 0 {
	t.FailNow()
}
```
