- Passwprd string `json:"-"` ทำแบบนี้ข้อมมูลไม่ออก fontend
-เอา struct ไม่ตรงประเภทแต่มีข้อมูลบาง ฟิวเหมือนกันจะดึงค่าจาก db ได้ไหม
- ถ้า verfily ไม่ทำเป็น middleware รูปแบบ func จะเป็นยังไง 
    -ทำ Func ที่ return Err ออกมา
    -จะใช้กับ middleware แค่ ถ้ามัน err ก็ไม่ Next

- create user ต้องไปเช็ค username ก่อนว่ามีคนใช้รึยัง  ***
- เทส ส่ง token ที่ user ไม่มี จะเกิดอะไรขึ้น ไปสรา้ง token ในเว็บ
    -ทำไม่ได้

- ดูเรื่อง claims  ของ youtube มันทำงานยังไง
    -youtube ใช้  parsewithclaims มันเลย  เซตค่าต่างกัน

- ดูเรื่อง time.time ดูคำสั่ง  time.pare()

- ย้าย validation เป็นไฟล์แยกออกไปใน util ทำแล้ว
    -การใช้ Validation 
        -ใช้โดยตรงผ่าน util
        -ใช้ผ่าน layer service ทำ interface ไว้ใน layer service **** 
    -ลองส่งค่าเป็น pointer ดูสิว่าได้ไหม ******************
        -หน้า handler ส่งได้ แต่ ตัว Func Vlidation ต้องไม่เป็น * เพราะมันเป็น interface{}
        -ทำ  *interfaace{} แบบนี้ หน้า handler ส่งยังไงก็ error
        - ปิด CreateCart_Handler, CreateCart_Handler,
            
- product Fuc New ต่างๆจะส่ง pointer เข้ามา มันต่างกันยังไง
    -ประหยัด ทรัพยากร เพราะ มันส่งค่า ที่อยู่ของตัวแปรนั้นไปแทน

- Validation Data ยังไม่เช็ค การกด space bar ***********

- pagination  ลองปรับแก้ให้เหมาะสม Search ใช้

- err ปั้นที่ handler ต้องย้ายไปที่service ***************
    -handler return แค่ err พอ เพราะทำหน้ารบและเช็คและ ปั้นข้อมูลส่ง service พอแล้ว
    
-gorm  find จะไม่มี  err กลับมา อย่างอื่นมี err / first take last
    -ใช้ order ช้วย เรียงลำดับ

- get cart for user
    -ดึง carts กับ cart_items แล้วเอาข้อมูลมา map กัน
    - เอา cart_items มาจัดกลุ่มก่อน
    - เอา cart_items มา map กับ carts
- get cart for store
    -ดึง DB ครั้งเดียว แต่ repository มีการดึง 2 ฐานข้อมูล
    - และเอาข้อมูลมา map กัน
     - เอา cart_items มาจัดกลุ่มก่อน
    - เอา cart_items มา map กับ carts

- ปรับตัวแปรให้ return pointer *****************

- แก้ response error cart, product, store เหลือ user *****************
    -validateDataUser เอา ชื่อ struct มาต่อ string น่าจะทำให้รู้ว่า error จาก ตรงไหน

-func เฉพาะด้านจะอยู่ใน  models  เช่น  ราคา promotion จะอยู่ใน models product *********
