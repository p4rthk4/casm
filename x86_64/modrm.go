// Copyright (c) 2023 Parth Degama
// This code is licensed under MIT license

// modrm

package x86_64

import "fmt"

// add modRM byte
func addModRM(opcode *archOpcode, inst *instruction, bitMode int, pf *prefix) ([]uint8, error) {

	// calculate modRM and return
	if len(inst.operands) != 2 {
		return nil, fmt.Errorf("internal error: todo: operand size is not two")
	}

	modRMrmOper := &operand{t: undefinedOperand}  // modrm r/m operand
	modRMregOper := &operand{t: undefinedOperand} // modrm reg operand

	// loop of arch operands
	for aOperIndex, aOper := range opcode.operands {
		if isMemoryOperand(aOper.t) {
			/*
				if arch operand is memory operand then
				assign inst operand to modRMrmOper
			*/
			modRMrmOper = &inst.operands[aOperIndex]
		}
		if isRegOperand(aOper.t) {
			/*
				if arch operand is register operand then
				assign inst operand to modRMregOper
			*/
			modRMregOper = &inst.operands[aOperIndex]
		}
	}

	// modRMreg info
	modRMregInfo, err := registerInfo(int(modRMregOper.v))
	if err != nil {
		return nil, err
	}

	// check modRMreg is valid or not
	err = registerIsValid(modRMregInfo, bitMode)
	if err != nil {
		return nil, err
	}

	// check operand override prefix
	err = checkOperandOverride(modRMregOper, bitMode, pf)
	if err != nil {
		return nil, err
	}

	// get reg field
	regField, err := modRMregField(*modRMregOper)
	if err != nil {
		return nil, err
	}

	fmt.Println(modRMrmOper, modRMregOper, regField)

	// calc modrm
	modRMByte, err := calcModRM(modRMrmOper, regField, bitMode, pf)
	if err != nil {
		return nil, err
	}

	return modRMByte, nil
}

// calc modrm
func calcModRM(rmOper *operand, regField int, bitMode int, pf *prefix) ([]uint8, error) {

	modrmBytes := []uint8{} // modrm bytes

	if rmOper.t == mem {
		// todo: modrm mem operand support
		//fmt.Println(rmOper.m)

		switch len(rmOper.m) {
		case 1:
			// if only one operand in memory operand
			memOper := rmOper.m[0]

			if isRegOperand(memOper.t) {

				// get mem register info
				regInfo, err := registerInfo(int(memOper.v))
				if err != nil {
					return nil, err
				}

				// check register is valid in mem or not
				err = registerIsValidInMem(regInfo, bitMode)
				if err != nil {
					return nil, err
				}

				// check address override prefix
				err = checkAddressOverride(&memOper, bitMode, pf)
				if err != nil {
					return nil, err
				}

				// get mem field
				modrmMemField, err := modRMmemRegField(regInfo)
				if err != nil {
					return nil, err
				}

				// gen modrm byte
				modRMByte := modRMbyte(0b00, regField, modrmMemField)
				modrmBytes = append(modrmBytes, uint8(modRMByte))

				return modrmBytes, nil // return modrm byte

			}

			return nil, fmt.Errorf("todo: modrm disp mem opernad not support")

		default:
			//
			return nil, fmt.Errorf("todo: modrm more then one mem opernad not support")
		}

	}

	if isRegOperand(rmOper.t) {
		// if rm operand is register

		// get r/m register info
		modRMrmRegInfo, err := registerInfo(int(rmOper.v))
		if err != nil {
			// if error then return error
			return nil, err
		}

		// check modRMreg is valid or not
		err = registerIsValid(modRMrmRegInfo, bitMode)
		if err != nil {
			return nil, err
		}

		// gen modrm byte
		modRMByte := modRMbyte(0b11, regField, modRMrmRegInfo.baseOffset)
		modrmBytes = append(modrmBytes, uint8(modRMByte))

		return modrmBytes, nil // return modrm byte

	}

	return nil, nil
}

// modrm reg field
func modRMregField(oper operand) (int, error) {

	if !isRegOperand(oper.t) {
		// if opernad is not register then return error
		return 0, fmt.Errorf("is not modrm reg type")
	}

	// get register info
	regInfo, err := registerInfo(int(oper.v))
	if err != nil {
		// if error then return error
		return 0, nil
	}

	/*
		return register base offset because
		register base offset is modrm reg field
	*/

	return regInfo.baseOffset, nil
}

// modRM mem field
func modRMmemRegField(r register) (int, error) {

	switch r.bitSize {
	case 16:

		// if register is 16 bit
		return 0, fmt.Errorf("todo: 16-bit register")

	case 32, 64:
		// if register is 32 or 64 bit

		switch r.baseOffset {
		case 4, 5:
			/*
				if register is esp, ebp, rsp or rbp
				this not support single reg mem
			*/
			return 0, fmt.Errorf("todo: %v register modrm byte", r.name)
		}

		return r.baseOffset, nil

	}

	return 0, fmt.Errorf("invalid modrm reg mem")

}

// get modrm byte
func modRMbyte(mod int, reg int, rm int) int {

	modRMByte := mod<<6 | reg<<3 | rm<<0

	return modRMByte

}
