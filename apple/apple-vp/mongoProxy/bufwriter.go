/*
 * @Author: fyfishie
 * @Date: 2023-03-21:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-15:15
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package mongoProxy

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"vp/lib"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BufferWriter struct {
	buffer []mongo.WriteModel
	bufLen int
	proxy  *Proxy
}

func NewBufferWriter(proxy *Proxy, bufLen int) *BufferWriter {
	return &BufferWriter{
		buffer: []mongo.WriteModel{},
		bufLen: bufLen,
		proxy:  proxy,
	}
}
func NewBufferWriterFromAddr(mongoAddr lib.MongoCollectionAddr, bufLen int) (*BufferWriter, error) {
	p := NewProxy(mongoAddr)
	err := p.Connect()
	if err != nil {
		return nil, err
	}
	return NewBufferWriter(p, bufLen), nil
}

// write model to mongodb with buffer
func (w *BufferWriter) WriteOne(model mongo.WriteModel) error {
	w.buffer = append(w.buffer, model)
	if len(w.buffer) >= w.bufLen {
		return w.Flush()
	}
	return nil
}

func (w *BufferWriter) Flush() error {
	opts := options.BulkWrite().SetOrdered(false)
	_, err := w.proxy.collection.BulkWrite(context.TODO(), w.buffer, opts)
	w.buffer = []mongo.WriteModel{}
	return err
}

func (b *BufferWriter) MvSpingRes2Mongo(respath string, UID string) error {
	rfi, err := os.OpenFile(respath, os.O_RDONLY, 0000)
	if err != nil {
		return err
	}
	defer rfi.Close()
	rdr := bufio.NewReader(rfi)
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		rsping := lib.RawSping{}
		if json.Unmarshal(line, &rsping) != nil {
			continue
		}
		data := lib.SpingWithUID{
			VPUID: UID,
			IP:    rsping.IP,
			Ttl:   strconv.Itoa(rsping.Ttl),
		}
		b.WriteOne(mongo.NewInsertOneModel().SetDocument(data))
	}
	return b.Flush()
}

// close client in mongo driver
func (b *BufferWriter) Close() error {
	return b.proxy.Disconnect()
}
