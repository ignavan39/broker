import { BeforeInsert, BeforeUpdate, Column, Entity } from "typeorm";
import * as crypto from 'crypto';

@Entity({name:'user'})
export class User {
    @Column({type:'uuid'})
    id: string

    @Column({type:'text'})
    email: string;

    @Column({type: 'text'})
    password: string;

    @BeforeInsert()
    hashPasswordBeforeInsert() {
      this.password = crypto.createHmac('sha256', this.password).digest('hex');
    }
  
    @BeforeUpdate()
    hashPasswordBeforeUpdate() {
      if (this.password) {
        this.password = crypto.createHmac('sha256', this.password).digest('hex');
      }
    }
}