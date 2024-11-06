import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { User, UserStatus } from '../user';

const testUsers = [
  {user_id: 1, user_name: "johndoe", first_name: "John", last_name: "Doe", email:"johndoe@gmail.com", user_status: UserStatus.Active, department:"IT"},
  {user_id: 2, user_name: "janedoe", first_name: "Jane", last_name: "Doe", email:"janedoe@gmail.com", user_status: UserStatus.Inactive, department:"IT"}
]

export type Response = {
    code:    number,
		message: string,
		details?: string,
		data?: User | User[] | number,
}

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = 'http://localhost:8080';
  private usersUrl = `${this.apiUrl}/api/v1/users`;

  constructor(private http: HttpClient) { }

  getUser(user_id: number): Observable<Response> {
    // return of(testUsers[0]);
    
    return this.http.get<Response>(`${this.usersUrl}/${user_id}`);
  }

  getUsers(): Observable<Response> {
    // return of(testUsers);
    
    return this.http.get<Response>(this.usersUrl);
  }

  createUser(user: User): Observable<Response> {
    return this.http.post<Response>(this.usersUrl, user);
  }

  deleteUser(user_id: number): Observable<Response> {
    return this.http.delete<Response>(`${this.usersUrl}/${user_id}`);
  }

  updateUser(user: User): Observable<Response> {
    return this.http.put<Response>(this.usersUrl, user);
  }
}