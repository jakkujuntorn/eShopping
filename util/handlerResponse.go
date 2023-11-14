package util

// ******** Success
type response_Data struct {
	Success string
	Data    interface{}
	JwtToken string
}

func Data_Response(data interface{}) response_Data {
	return response_Data{
		Success: "Success",
		Data:data,
	}
}

// ถ้า Error ให้ return pattern นี้ออกไป อยู่ที่ handler_Error
// Success:"false"
// Error: {...}


// มีข้อมูลแนบไปด้วย
// "Success": "true",
// "Data": dataUser,
// "jwtToken":"" // ถ้ามี

// ส่งแค่ message
//"success":"true",
//"Message":"Create User Success",




