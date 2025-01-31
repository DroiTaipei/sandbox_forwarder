// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package sampling

import (
	"bytes"
	"fmt"
	"github.com/DroiTaipei/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type SamplingManager interface {
	// Parameters:
	//  - ServiceName
	GetSamplingStrategy(serviceName string) (r *SamplingStrategyResponse, err error)
}

type SamplingManagerClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewSamplingManagerClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *SamplingManagerClient {
	return &SamplingManagerClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewSamplingManagerClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *SamplingManagerClient {
	return &SamplingManagerClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - ServiceName
func (p *SamplingManagerClient) GetSamplingStrategy(serviceName string) (r *SamplingStrategyResponse, err error) {
	if err = p.sendGetSamplingStrategy(serviceName); err != nil {
		return
	}
	return p.recvGetSamplingStrategy()
}

func (p *SamplingManagerClient) sendGetSamplingStrategy(serviceName string) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("getSamplingStrategy", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := SamplingManagerGetSamplingStrategyArgs{
		ServiceName: serviceName,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *SamplingManagerClient) recvGetSamplingStrategy() (value *SamplingStrategyResponse, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "getSamplingStrategy" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "getSamplingStrategy failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "getSamplingStrategy failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error1 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error2 error
		error2, err = error1.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error2
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "getSamplingStrategy failed: invalid message type")
		return
	}
	result := SamplingManagerGetSamplingStrategyResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type SamplingManagerProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      SamplingManager
}

func (p *SamplingManagerProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *SamplingManagerProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *SamplingManagerProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewSamplingManagerProcessor(handler SamplingManager) *SamplingManagerProcessor {

	self3 := &SamplingManagerProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self3.processorMap["getSamplingStrategy"] = &samplingManagerProcessorGetSamplingStrategy{handler: handler}
	return self3
}

func (p *SamplingManagerProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x4 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x4.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x4

}

type samplingManagerProcessorGetSamplingStrategy struct {
	handler SamplingManager
}

func (p *samplingManagerProcessorGetSamplingStrategy) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := SamplingManagerGetSamplingStrategyArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getSamplingStrategy", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := SamplingManagerGetSamplingStrategyResult{}
	var retval *SamplingStrategyResponse
	var err2 error
	if retval, err2 = p.handler.GetSamplingStrategy(args.ServiceName); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getSamplingStrategy: "+err2.Error())
		oprot.WriteMessageBegin("getSamplingStrategy", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("getSamplingStrategy", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - ServiceName
type SamplingManagerGetSamplingStrategyArgs struct {
	ServiceName string `thrift:"serviceName,1" json:"serviceName"`
}

func NewSamplingManagerGetSamplingStrategyArgs() *SamplingManagerGetSamplingStrategyArgs {
	return &SamplingManagerGetSamplingStrategyArgs{}
}

func (p *SamplingManagerGetSamplingStrategyArgs) GetServiceName() string {
	return p.ServiceName
}
func (p *SamplingManagerGetSamplingStrategyArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *SamplingManagerGetSamplingStrategyArgs) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.ServiceName = v
	}
	return nil
}

func (p *SamplingManagerGetSamplingStrategyArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("getSamplingStrategy_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *SamplingManagerGetSamplingStrategyArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("serviceName", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:serviceName: ", p), err)
	}
	if err := oprot.WriteString(string(p.ServiceName)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.serviceName (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:serviceName: ", p), err)
	}
	return err
}

func (p *SamplingManagerGetSamplingStrategyArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SamplingManagerGetSamplingStrategyArgs(%+v)", *p)
}

// Attributes:
//  - Success
type SamplingManagerGetSamplingStrategyResult struct {
	Success *SamplingStrategyResponse `thrift:"success,0" json:"success,omitempty"`
}

func NewSamplingManagerGetSamplingStrategyResult() *SamplingManagerGetSamplingStrategyResult {
	return &SamplingManagerGetSamplingStrategyResult{}
}

var SamplingManagerGetSamplingStrategyResult_Success_DEFAULT *SamplingStrategyResponse

func (p *SamplingManagerGetSamplingStrategyResult) GetSuccess() *SamplingStrategyResponse {
	if !p.IsSetSuccess() {
		return SamplingManagerGetSamplingStrategyResult_Success_DEFAULT
	}
	return p.Success
}
func (p *SamplingManagerGetSamplingStrategyResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *SamplingManagerGetSamplingStrategyResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *SamplingManagerGetSamplingStrategyResult) readField0(iprot thrift.TProtocol) error {
	p.Success = &SamplingStrategyResponse{}
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *SamplingManagerGetSamplingStrategyResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("getSamplingStrategy_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *SamplingManagerGetSamplingStrategyResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *SamplingManagerGetSamplingStrategyResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SamplingManagerGetSamplingStrategyResult(%+v)", *p)
}
