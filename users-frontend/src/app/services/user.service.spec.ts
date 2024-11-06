import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { UserService } from './user.service';
import { UserStatus } from '../user';
import { firstValueFrom } from 'rxjs';
import { httpTestWrap } from '../helpers/tests';

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

    const req = httpTestingController.expectOne(`${host}/api/v1/users`);
    expect(req.request.method).toBe('GET');
    
    req.flush(httpTestWrap(testUsers));

    expect((await usersPromise).data).toEqual(testUsers);

    httpTestingController.verify();
  });

  it('should retrieve user from the API', async () => {
    const users$ = service.getUser(1);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/v1/users/1`);
    expect(req.request.method).toBe('GET');
    
    req.flush(httpTestWrap(testUsers[0]));

    expect((await usersPromise).data).toEqual(testUsers[0]);

    httpTestingController.verify();
  });

  it('should create user using API', async () => {
    const users$ = service.createUser(testUsers[0]);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/v1/users`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toBe(testUsers[0])
    
    req.flush(httpTestWrap(testUsers[0]));

    expect((await usersPromise).data).toEqual(testUsers[0]);

    httpTestingController.verify();
  });

  it('should update user using API', async () => {
    const users$ = service.updateUser(testUsers[0]);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/v1/users`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toBe(testUsers[0])
    
    req.flush(httpTestWrap(testUsers[0]));

    expect((await usersPromise).data).toBe(testUsers[0]);

    httpTestingController.verify();
  });

  it('should delete user using API', async () => {
    const users$ = service.deleteUser(testUsers[0].user_id);
    const usersPromise = firstValueFrom(users$);

    const req = httpTestingController.expectOne(`${host}/api/v1/users/${testUsers[0].user_id}`);
    expect(req.request.method).toBe('DELETE');
    
    req.flush(httpTestWrap(testUsers[0]));

    await usersPromise;

    httpTestingController.verify();
  });
});
