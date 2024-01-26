-จำลองเว็บไซด์ขายของออนไลด์ 
    -Gin
    -Hexagonal Architecture
        -port and adaptors
    -DB
        -mysql (user, product, store)
        -Postgres (cart)
        -Redis (product)
    -Code Structure
        -By Layer
    
********************************
-สิ่งที่จะทำต่อ
     -API เกี่ยวกับ create อยากเขียนแบบใหม่ที่ไม่ใช้ เช็คข้อมูลซ้ำแล้ว create ลงไปเลย (ตอนนี้ยังคิดไม่ออก)
    -refactor code 
    -ทำ unit test 
    -แก้ logic ที่ cart service (ให้ for ชั้นเดียว)
    -ใช้ kafka ตรง register หลังจาก register เสร็จให้ส่งเมล

*******************************
- Note
    -โปรเจคนี้เป็นการจำลองการทำงาน และใส logic ที่แปลกๆลงไป เพื่อฝึกการ query, func บาง func ถูกเอาไปใสในงานจริง
    -DB ไม่ทำ foreign key(ตั้งใจทำแบบนี้) เพื่อการ Query ที่ยากขึ้น เพราะต้องเอาข้อมูลมา map กันเองภายนอก
    -ในโปรเจคนี้จะมีการ comment code ที่เยอะ เพราะเอาไว้อ่านตอนมาแก้จะได้เข้าใจได้ง่าย 
    -ถ้าใครได้เข้ามาดูฝาก รีวิว logic บาง func ด้วยว่ามันควรจะเขียนยังไงให้ดีกว่านี้  ************
        -โดยเฉพาะ "repository/cart_postgres.go" (มีการทำ Transaction)
