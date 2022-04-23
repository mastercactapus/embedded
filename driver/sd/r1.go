package sd

type R1 byte

func (r R1) IsParamErr() bool      { return r&(1<<6) != 0 }
func (r R1) IsAddrErr() bool       { return r&(1<<5) != 0 }
func (r R1) IsEraseSeqErr() bool   { return r&(1<<4) != 0 }
func (r R1) IsCRCErr() bool        { return r&(1<<3) != 0 }
func (r R1) IsIllegalCmdErr() bool { return r&(1<<2) != 0 }
func (r R1) IsEraseReset() bool    { return r&(1<<1) != 0 }
func (r R1) IsIdleSet() bool       { return r&(1<<0) != 0 }
