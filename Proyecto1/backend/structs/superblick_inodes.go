package structs

import "encoding/binary"

type Inodes struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [16]int32
	I_type  [1]byte
	I_perm  int32
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type FolderBlock struct {
	B_content [4]Content
}

type Fileblock struct {
	B_content [64]byte
}

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [19]byte
	S_umtime            [19]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

func NewSuperBlock() SuperBlock {
	return SuperBlock{
		S_filesystem_type:   0,
		S_inodes_count:      0,
		S_blocks_count:      0,
		S_free_blocks_count: 0,
		S_free_inodes_count: 0,
		S_mtime:             [19]byte{},
		S_umtime:            [19]byte{},
		S_mnt_count:         0,
		S_magic:             0xEF53,
		S_inode_size:        int32(binary.Size(Inodes{})),
		S_block_size:        int32(binary.Size(FolderBlock{})),
		S_first_ino:         0,
		S_first_blo:         0,
		S_bm_inode_start:    0,
		S_bm_block_start:    0,
		S_inode_start:       0,
		S_block_start:       0,
	}
}

func NewInodes() Inodes {
	return Inodes{
		I_uid:   -1,
		I_gid:   -1,
		I_size:  -1,
		I_atime: [19]byte{},
		I_ctime: [19]byte{},
		I_mtime: [19]byte{},
		I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'-'},
		I_perm:  -1,
	}
}
