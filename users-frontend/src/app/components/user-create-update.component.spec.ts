// Edit page loads data
// Create page happy path

import { ComponentFixture, TestBed, tick } from "@angular/core/testing";
import { UserCreateUpdateComponent } from "./user-create-update.component";
import { UserService } from "../services/user.service";
import { ActivatedRoute, Router } from "@angular/router";
import { MatButtonModule } from "@angular/material/button";
import { NO_ERRORS_SCHEMA } from "@angular/core";
import { of } from "rxjs";
import { UserStatus } from "../user";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { By } from "@angular/platform-browser";
import { provideHttpClient } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { MatSelectHarness } from "@angular/material/select/testing";
import { HarnessLoader } from "@angular/cdk/testing";
import { MatButtonHarness } from "@angular/material/button/testing";
import { MatInputHarness } from "@angular/material/input/testing";

// Create page sad path
const testUsers = [
  {user_id: 1, user_name: "johndoe", first_name: "John", last_name: "Doe", email:"johndoe@gmail.com", user_status: UserStatus.Active, department:"IT"},
  {user_id: 2, user_name: "janedoe", first_name: "Jane", last_name: "Doe", email:"janedoe@gmail.com", user_status: UserStatus.Inactive, department:"IT"}
]

// describe('UserCreateUpdateComponent - Create', () => {

// });

describe('UserCreateUpdateComponent - Create', () => {
    let component: UserCreateUpdateComponent;
    let fixture: ComponentFixture<UserCreateUpdateComponent>;
    let userService: jasmine.SpyObj<UserService>;
    let routerSpy: jasmine.SpyObj<Router>;
    let routeSpy: jasmine.SpyObj<ActivatedRoute>;
    let httpTestingController: HttpTestingController;
    let loader: HarnessLoader;
  
    beforeEach(async () => {
      routerSpy = jasmine.createSpyObj('Router', ['navigate'], { url: '/create' });
      routeSpy = jasmine.createSpyObj('ActivatedRoute', ['queryParams']);
      
      await TestBed.configureTestingModule({
        imports: [UserCreateUpdateComponent, MatButtonModule, NoopAnimationsModule],
        providers: [
          { provide: Router, useValue: routerSpy },
          { provide: ActivatedRoute, useValue: { queryParams: of({})}},
          provideHttpClient(), 
          provideHttpClientTesting()
        ],
        schemas: [NO_ERRORS_SCHEMA]
      }).compileComponents();
  
      fixture = TestBed.createComponent(UserCreateUpdateComponent);
      component = fixture.componentInstance;
      userService = TestBed.inject(UserService) as jasmine.SpyObj<UserService>;
      httpTestingController = TestBed.inject(HttpTestingController);
      loader = TestbedHarnessEnvironment.loader(fixture);
      
      fixture.detectChanges();
    });
  
    afterEach(() => {
      httpTestingController.verify();
      fixture.destroy(); fixture.destroy();
    });
  
    it('should fetch user data when in edit mode and populate inputs', async () => {
        const submitButton = await loader.getHarness(MatButtonHarness.with({text: 'Submit'}));
        const inputs = await loader.getAllHarnesses(MatInputHarness);
        const select = await loader.getHarness(MatSelectHarness);
        expect(inputs.length).toBe(5);
        expect(await submitButton.isDisabled()).toBeTrue();

        await inputs[0].setValue('johndoe');
        await inputs[1].setValue('John');
        await inputs[2].setValue('Doe');
        await inputs[3].setValue('johndoe@gmail.com');
        await inputs[4].setValue('IT');

        await select.open();
        const bugOption = await select.getOptions({text: 'Active'});
        await bugOption[0].click();

        fixture.detectChanges();
        expect(await submitButton.isDisabled()).toBeFalse();

        await submitButton.click();

        await fixture.whenStable(); // Waits for all async operations to finish

        // Detect changes to update the view
        fixture.detectChanges();

        const req = httpTestingController.expectOne(`http://localhost:8080/api/v1/users`);
        expect(req.request.method).toBe('POST');
        req.flush(testUsers[0]);

        httpTestingController.verify();
    });
  });

describe('UserCreateUpdateComponent - Edit', () => {
  let component: UserCreateUpdateComponent;
  let fixture: ComponentFixture<UserCreateUpdateComponent>;
  let userService: jasmine.SpyObj<UserService>;
  let routerSpy: jasmine.SpyObj<Router>;
  let routeSpy: jasmine.SpyObj<ActivatedRoute>;
  let httpTestingController: HttpTestingController;
  let loader: HarnessLoader;

  beforeEach(async () => {
    routerSpy = jasmine.createSpyObj('Router', ['navigate'], { url: '/edit' });
    routeSpy = jasmine.createSpyObj('ActivatedRoute', ['queryParams']);
    
    await TestBed.configureTestingModule({
      imports: [UserCreateUpdateComponent, MatButtonModule, NoopAnimationsModule],
      providers: [
        { provide: Router, useValue: routerSpy },
        { provide: ActivatedRoute, useValue: { queryParams: of({'user_id': '1'})}},
        provideHttpClient(), 
        provideHttpClientTesting()
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();

    fixture = TestBed.createComponent(UserCreateUpdateComponent);
    component = fixture.componentInstance;
    userService = TestBed.inject(UserService) as jasmine.SpyObj<UserService>;
    httpTestingController = TestBed.inject(HttpTestingController);
    loader = TestbedHarnessEnvironment.loader(fixture);
    
    fixture.detectChanges();
  });

  afterEach(() => {
    httpTestingController.verify();
    fixture.destroy();
  });

  it('should fetch user data when in edit mode and populate inputs', async () => {
    fixture.detectChanges();

    const req = httpTestingController.expectOne('http://localhost:8080/api/v1/users/1');
    expect(req.request.method).toBe('GET'); 
    req.flush({data: testUsers[0]}); 

    expect(component.user).toEqual(testUsers[0]);
    fixture.detectChanges();

    const inputs = await loader.getAllHarnesses(MatInputHarness);
    const select = await loader.getHarness(MatSelectHarness);

    expect(await inputs[0].getValue()).toBe('johndoe');
    expect(await inputs[1].getValue()).toBe('John');
    expect(await inputs[2].getValue()).toBe('Doe');
    expect(await inputs[3].getValue()).toBe('johndoe@gmail.com');
    expect(await inputs[4].getValue()).toBe('IT');

    expect(await select.getValueText()).toBe('Active');
  });
});
