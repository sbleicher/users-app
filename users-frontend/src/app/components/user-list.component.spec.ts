import { ComponentFixture, TestBed } from '@angular/core/testing';
import { UserListComponent } from './users-list.component';
import { UserService } from '../services/user.service';
import { of } from 'rxjs';
import { User, UserStatus, userStatusToString } from '../user';
import { By } from '@angular/platform-browser';
import { DebugElement, NO_ERRORS_SCHEMA } from '@angular/core';
import { Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';

const HEADER_ROW = 1
const testUsers = [
  {user_id: 1, user_name: "johndoe", first_name: "John", last_name: "Doe", email:"johndoe@gmail.com", user_status: UserStatus.Active, department:"IT"},
  {user_id: 2, user_name: "janedoe", first_name: "Jane", last_name: "Doe", email:"janedoe@gmail.com", user_status: UserStatus.Inactive, department:"IT"}
]

describe('UserListComponent', () => {
  let component: UserListComponent;
  let fixture: ComponentFixture<UserListComponent>;
  let userService: jasmine.SpyObj<UserService>;
  let routerSpy: jasmine.SpyObj<Router>;

  beforeEach(async () => {
    const userServiceMock = jasmine.createSpyObj('UserService', ['getUsers', 'deleteUser']);
    routerSpy = jasmine.createSpyObj('Router', ['navigate']);
    
    await TestBed.configureTestingModule({
      imports: [UserListComponent, MatButtonModule],
      providers: [
        { provide: UserService, useValue: userServiceMock },
        { provide: Router, useValue: routerSpy }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();

    fixture = TestBed.createComponent(UserListComponent);
    component = fixture.componentInstance;
    userService = TestBed.inject(UserService) as jasmine.SpyObj<UserService>;

    userService.getUsers.and.returnValue(of(testUsers));
    userService.deleteUser.and.returnValue(of());
    
    fixture.detectChanges();
  });

  it('should fetch data from the service and display it in table', () => {
    expect(userService.getUsers).toHaveBeenCalled(); 

    const rows = fixture.debugElement.queryAll(By.css('tr'));
    expect(rows.length).toBe(HEADER_ROW + testUsers.length);
    
    testRow(rows[1].queryAll(By.css('td')), testUsers[0]);
    testRow(rows[2].queryAll(By.css('td')), testUsers[1]);
  });

  it('should navigate to the create page', () => {
    const buttons = fixture.debugElement.queryAll(By.css('button'));
    buttons[0].triggerEventHandler('click', null);

    expect(routerSpy.navigate).toHaveBeenCalledWith(['/create']);
  });
  
  it('should navigate to the edit page with user_id', () => {
    const buttons = fixture.debugElement.queryAll(By.css('button'));
    buttons[1].triggerEventHandler('click', null);

    expect(routerSpy.navigate).toHaveBeenCalledWith(['/edit'], { queryParams: { user_id: testUsers[0].user_id } });
  });

  it('should call the delete api endpoint and refresh page', () => {
    const buttons = fixture.debugElement.queryAll(By.css('button'));
    buttons[2].triggerEventHandler('click', null);

    expect(userService.deleteUser).toHaveBeenCalledWith(testUsers[0].user_id); 
  });

  function testRow(row: DebugElement[], user: User){
    expect(row[0].nativeElement.textContent.trim()).toBe(String(user.user_id));
    expect(row[1].nativeElement.textContent.trim()).toBe(user.user_name);
    expect(row[2].nativeElement.textContent.trim()).toBe(user.first_name);
    expect(row[3].nativeElement.textContent.trim()).toBe(user.last_name);
    expect(row[4].nativeElement.textContent.trim()).toBe(user.email);
    expect(row[5].nativeElement.textContent.trim()).toBe(userStatusToString(user.user_status));
    expect(row[6].nativeElement.textContent.trim()).toBe(user.department);
  }
});
