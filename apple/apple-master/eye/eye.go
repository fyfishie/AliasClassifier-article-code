/*
 * @Author: fyfishie
 * @Date: 2023-03-27:18
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-03-27:18
 * @Description: monitor status of task, in application layer
 * @email: fyfishie@outlook.com
 */
package eye

type Eye[DataType any] struct {
	payload DataType
}

func (e *Eye[T]) Update(Data T) {
	e.payload = Data
}

func (e *Eye[T]) Read() (Data T) {
	return e.payload
}

func (e *Eye[T]) GetHandle() (handle *T) {
	return &e.payload
}
