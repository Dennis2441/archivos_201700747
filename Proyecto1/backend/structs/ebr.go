package structs

type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_size  int32
	Part_next  int32
	Part_name  [16]byte
}

func NewEBR() EBR {
	return EBR{
		Part_mount: [1]byte{'0'},
		Part_fit:   [1]byte{'w'},
		Part_start: -1,
		Part_size:  0,
		Part_next:  -1,
		Part_name:  [16]byte{'~', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}
