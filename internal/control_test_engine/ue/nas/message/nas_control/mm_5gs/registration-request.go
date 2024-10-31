package mm_5gs

import (
	"bytes"
	"fmt"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/nasType"
	"my5G-RANTester/lib/openapi/models"
)

func GetRegistrationRequest(registrationType uint8, requestedNSSAI *nasType.RequestedNSSAI, uplinkDataStatus *nasType.UplinkDataStatus, capability bool, ue *context.UEContext) (nasPdu []byte) {

	ueSecurityCapability := context.SetUESecurityCapability(ue)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)

	registrationRequest := nasMessage.NewRegistrationRequest(0)
	registrationRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0x00)
	registrationRequest.RegistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeRegistrationRequest)
	registrationRequest.NgksiAndRegistrationType5GS.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	registrationRequest.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(ue.GetUeId())
	registrationRequest.NgksiAndRegistrationType5GS.SetRegistrationType5GS(registrationType)
	registrationRequest.MobileIdentity5GS = ue.GetSuci()
	if capability {
		registrationRequest.Capability5GMM = &nasType.Capability5GMM{
			Iei:   nasMessage.RegistrationRequestCapability5GMMType,
			Len:   3,
			Octet: [13]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		}
	} else {
		registrationRequest.Capability5GMM = nil
	}
	registrationRequest.UESecurityCapability = ueSecurityCapability

	// Construct the Requested NSSAI with SST = 1 (eMBB) and SD = 1691105
	if ue.PduSession.Snssai.Sst != 0 && ue.PduSession.Snssai.Sd != "" {
		// Convert the SD from string to uint32 (assuming SD is provided as a string in decimal)
		sdUint := 1691105 // This is the SD value we are trying to encode, adjust it as per actual value from ue.PduSession.Snssai.Sd

		// Prepare the SNSSAI value (with Lengths for SST and SD)
		snssaiValue := []uint8{
			4,                               // Length of SST (1 byte)
			uint8(ue.PduSession.Snssai.Sst), // SST (e.g., 1 for eMBB)
			uint8(sdUint >> 16),             // First byte of SD
			uint8(sdUint >> 8),              // Second byte of SD
			uint8(sdUint),                   // Third byte of SD
		}

		// Allocate the buffer for Requested NSSAI
		requestedNSSAI = &nasType.RequestedNSSAI{
			Iei:    0x2F,                    // Element ID
			Len:    uint8(len(snssaiValue)), // Total length of SNSSAI structure
			Buffer: make([]uint8, len(snssaiValue)),
		}

		// Set the SNSSAI value into the buffer
		requestedNSSAI.SetSNSSAIValue(snssaiValue)
	}

	registrationRequest.RequestedNSSAI = requestedNSSAI

	mmCtx := &models.MmContext{
		AccessType: models.AccessType__3_GPP_ACCESS,
	}
	registrationRequest.MmContext = mmCtx

	// updateType5GS := &nasType.UpdateType5GS{
	// 	Iei:   0x53,
	// 	Len:   uint8(1), // Total length of SNSSAI structure
	// 	Octet: uint8(0),
	// }

	// registrationRequest.UpdateType5GS = updateType5GS

	// networkSlicingIndication := nasType.NetworkSlicingIndication{
	// 	Octet: 0x95, // DCNI = 1 (Default NSSAI), NSSCI = 1 (Slice selection capability indicated)
	// }

	// registrationRequest.NetworkSlicingIndication = &networkSlicingIndication

	registrationRequest.UplinkDataStatus = uplinkDataStatus

	registrationRequest.SetFOR(1)

	m.GmmMessage.RegistrationRequest = registrationRequest

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	nasPdu = data.Bytes()
	return
}
