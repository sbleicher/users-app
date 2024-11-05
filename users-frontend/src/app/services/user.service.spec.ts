import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { UserService } from './user.service';
import { UserStatus } from '../user';
import { firstValueFrom } from 'rxjs';

const testUsers = [
    {user_id: 1, user_name: "johndoe", first_name: "John", last_name: "Doe", email:"johndoe@gmail.com", user_status: UserStatus.Active, department:"IT"},
    {user_id: 2, user_name: "janedoe", first_name: "Jane", last_name: "Doe", email:"janedoe@gmail.com", user_status: UserStatus.Inactive, department:"IT"}
  ]
const host = "http://localhost:8080";

describe('UserService', () => {
  let service: UserService;
  let httpTestingController: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [UserService, provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(UserService);
    httpTestingController = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it('should retrieve users from the API', async () => {
    const users$ = service.getUsers();
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/users`);
    expect(req.request.method).toBe('GET');
    
    req.flush(testUsers);

    expect(await usersPromise).toEqual(testUsers);

    httpTestingController.verify();
  });

  it('should retrieve user from the API', async () => {
    const users$ = service.getUser(1);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/users/1`);
    expect(req.request.method).toBe('GET');
    
    req.flush(testUsers[0]);

    expect(await usersPromise).toEqual(testUsers[0]);

    httpTestingController.verify();
  });

  it('should create user using API', async () => {
    const users$ = service.createUser(testUsers[0]);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/users`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toBe(testUsers[0])
    
    req.flush(testUsers[0]);

    expect(await usersPromise).toEqual(testUsers[0]);

    httpTestingController.verify();
  });

  it('should update user using API', async () => {
    const users$ = service.updateUser(testUsers[0]);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/users/1`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toBe(testUsers[0])
    
    req.flush(testUsers[0]);

    expect(await usersPromise).toEqual(testUsers[0]);

    httpTestingController.verify();
  });

  it('should delete user using API', async () => {
    const users$ = service.deleteUser(testUsers[0].user_id);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/users/${testUsers[0].user_id}`);
    expect(req.request.method).toBe('DELETE');
    
    req.flush(testUsers[0]);

    await usersPromise;

    httpTestingController.verify();
  });

//   it('should handle a 404 error from the API', () => {
//     const errorMessage = 'Not Found';

//     // Call the service method that makes the HTTP request
//     service.getUsers().subscribe(
//       () => fail('should have failed with 404 error'),
//       (error) => {
//         expect(error.status).toBe(404);
//         expect(error.error).toBe(errorMessage);
//       }
//     );

//     // Expect that an HTTP GET request was made to '/api/users'
//     const req = httpTestingController.expectOne('/api/users');
    
//     // Respond with a 404 error
//     req.flush(errorMessage, { status: 404, statusText: 'Not Found' });

//     // Verify no other requests are pending
//     httpTestingController.verify();
//   });
});
