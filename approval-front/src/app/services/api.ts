import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080';

  private getAuthHeaders() {
    const token = localStorage.getItem('token');
    return new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });
  }

  // --- Public Endpoints ---

  login(credentials: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials);
  }

  createRequest(title: string): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/user/request`,
      { title },
      { headers: this.getAuthHeaders() },
    );
  }

  getMyRequests(): Observable<any> {
    return this.http.get(`${this.apiUrl}/user/my-requests`, { headers: this.getAuthHeaders() });
  }

  getAllRequests(): Observable<any> {
    return this.http.get(`${this.apiUrl}/admin/all-requests`, { headers: this.getAuthHeaders() });
  }

  approveRequests(ids: number[], status: string, reason: string): Observable<any> {
    const body = {
      ids: ids,
      status: status,
      admin_reason: reason,
    };
    return this.http.put(`${this.apiUrl}/admin/approve-multiple`, body, {
      headers: this.getAuthHeaders(),
    });
  }

  createUser(userData: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/admin/create-user`, userData, {
      headers: this.getAuthHeaders(),
    });
  }
}
