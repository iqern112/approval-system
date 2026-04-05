import { Component, inject, OnInit } from '@angular/core'; 
import { FormsModule } from '@angular/forms';
import { ApiService } from '../services/api'; 
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login implements OnInit {
  loginData = {
    username: '',
    password: '',
  };

  isLoading = false; 

  private apiService = inject(ApiService);
  private router = inject(Router);

  ngOnInit() {
    localStorage.clear();
  }

  onLogin() {
    if (!this.loginData.username || !this.loginData.password) {
      alert('กรุณากรอกข้อมูลให้ครบถ้วน');
      return;
    }

    this.isLoading = true; 

    this.apiService.login(this.loginData).subscribe({
      next: (res: any) => {
        // เก็บข้อมูลลงเครื่อง
        localStorage.setItem('token', res.token);
        localStorage.setItem('role', res.role);
        localStorage.setItem('username', res.username);

        this.loginData = { username: '', password: '' };

        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        console.error('Login Error:', err);
        alert('ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง');
        this.isLoading = false; 
      },
    });
  }
}