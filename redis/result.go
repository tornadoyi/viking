package redis



type Result struct {
	reply		interface{}
	error		error
}

func (h *Result) Reply() interface{} { return h.reply}

func (h *Result) Error() error { return h.error}

func (h *Result) Bool () (bool, error) { return Bool(h.reply, h.error) }
func (h *Result) ByteSlices () ([][]byte, error) { return ByteSlices(h.reply, h.error) }
func (h *Result) Bytes () ([]byte, error) { return Bytes(h.reply, h.error) }
func (h *Result) Float64 () (float64, error) { return Float64(h.reply, h.error) }
func (h *Result) Float64s () ([]float64, error) { return Float64s(h.reply, h.error) }
func (h *Result) Int () (int, error) { return Int(h.reply, h.error) }
func (h *Result) Int64 () (int64, error) { return Int64(h.reply, h.error) }
func (h *Result) Int64Map () (map[string]int64, error) { return Int64Map(h.reply, h.error) }
func (h *Result) Int64s () ([]int64, error) { return Int64s(h.reply, h.error) }
func (h *Result) IntMap () (map[string]int, error) { return IntMap(h.reply, h.error) }
func (h *Result) Ints () ([]int, error) { return Ints(h.reply, h.error) }
func (h *Result) MultiBulk () ([]interface{}, error) { return MultiBulk(h.reply, h.error) }
func (h *Result) Positions () ([]*[2]float64, error) { return Positions(h.reply, h.error) }
func (h *Result) String () (string, error) { return String(h.reply, h.error) }
func (h *Result) StringMap () (map[string]string, error) { return StringMap(h.reply, h.error) }
func (h *Result) Strings () ([]string, error) { return Strings(h.reply, h.error) }
func (h *Result) Uint64 () (uint64, error) { return Uint64(h.reply, h.error) }
func (h *Result) Values () ([]interface{}, error) { return Values(h.reply, h.error) }
func (h *Result) Interface () (interface{}, error) { return h.reply, h.error }
func (h *Result) IsNil () bool { return h.reply == nil }