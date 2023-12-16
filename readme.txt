-จำลองเว็บไซด์ขายของออนไลด์
    -Gin
    -Hexagonal Architecture
        -port and adaptors
    -DB
        -mysql
        -Postgres
        -Redis 

******************************
-สิ่งที่จะทำต่อ
    -ทำ unit test 
    -แก้ logic ที่ cart service (ให้ for ชั้นเดียว)
    -delete ยังไม่ได้ทำ มีแนวคิดแบบอื่น
        -แนวคิด จะลบจริงใน carts และ cart_items และจะเอาข้อมูลที่ลบไปใสใน ฐานข้อมูลใหม่ (delete_carts, delete_cart+items)
        -ทำ transaction (ใช้เทคนิคใน youtube  ทำแค่ func เดียว)
    -ใช้ kafka ตรง register หลังจาก register เสร็จให้ส่งเมล