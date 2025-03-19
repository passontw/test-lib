package basesk

type PacketHeader struct { // PacketHeader 定义
	Cmd  uint32 // 命令字
	Size uint32 // 数据包大小
	Seq  uint32 // 数据包序列号
}

func (p *PacketHeader) Init(cmd uint32, size uint32, seq uint32) {
	p.Cmd = cmd
	p.Size = size
	p.Seq = seq
}
