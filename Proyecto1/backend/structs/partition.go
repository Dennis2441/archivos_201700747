package structs

type Partition struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_size        int32
	Part_name        [16]byte
	Part_id          [4]byte
	Part_correlative int32
}

func NewPartition() Partition {
	return Partition{
		Part_status:      [1]byte{'0'},
		Part_type:        [1]byte{'p'},
		Part_fit:         [1]byte{'w'},
		Part_start:       -1,
		Part_size:        -1,
		Part_name:        [16]byte{'~', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Part_id:          [4]byte{'~', 0, 0, 0},
		Part_correlative: -1,
	}
}
