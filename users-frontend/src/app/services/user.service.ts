import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { User, UserStatus } from '../user';

const testUsers = [
  {user_id: 1, user_name: "johndoe", first_name: "John", last_name: "Doe", email:"johndoe@gmail.com", user_status: UserStatus.Active, department:"IT"},
  {user_id: 2, user_name: "janedoe", first_name: "Jane", last_name: "Doe", email:"janedoe@gmail.com", user_status: UserStatus.Inactive, department:"IT"}
]

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) { }

  getUser(user_id: number): Observable<User> {
    // return of(testUsers[0]);
    
    return this.http.get<User>(`${this.apiUrl}/api/users/${user_id}`);
  }

  getUsers(): Observable<User[]> {
    // return of(testUsers);
    
    return this.http.get<User[]>(`${this.apiUrl}/api/users`);
  }

  createUser(user: User): Observable<User> {
    return this.http.post<User>(`${this.apiUrl}/api/users`, user);
  }

  deleteUser(user_id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/api/users/${user_id}`);
  }

  updateUser(user: User): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/api/users/${user.user_id}`, user);
  }
}