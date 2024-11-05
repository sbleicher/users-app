import { Component } from '@angular/core';
import {User} from '../user';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatSelectModule} from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { UserService } from '../services/user.service';
import { ActivatedRoute, Router } from '@angular/router';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import { NgIf, NgClass  } from '@angular/common';

@Component({
    standalone: true,
    selector: 'user-create-update',
    templateUrl: './user-create-update.component.html',
    imports: [NgClass, FormsModule, MatFormFieldModule, MatInputModule, ReactiveFormsModule, MatGridListModule, MatSelectModule, MatButtonModule, NgIf, MatProgressSpinnerModule],
})
export class UserCreateUpdateComponent {
    form: FormGroup;

    user_id : number | undefined;
    user: User | undefined;
    mode: string = "create";

    constructor(private formBuilder: FormBuilder, private userService: UserService, private router: Router, private route: ActivatedRoute) {
        this.form = this.formBuilder.group({
            user_name: ['', Validators.required],
            first_name: ['', Validators.required],
            last_name: ['', Validators.required],
            email: ['', Validators.compose([Validators.required, Validators.email])],
            user_status: ['', Validators.required],
            department: [''],
          });

          if(this.router.url.includes("edit")){
            this.mode = "edit";
            this.route.queryParams.subscribe(params => {
                this.user_id = params['user_id'];
                if(this.user_id !== undefined){
                    this.getUser(this.user_id);
                } else {
                    this.router.navigate(['/']);
                }
            });
        }
    }

    navigateToListPage() {
        this.router.navigate(['/']);
    }

    getUser(user_id: number) {
        this.userService.getUser(user_id).subscribe({
            next: (data: any) => {
                this.user = data;
                this.updateFormWithUser(data);
            },
            error: (error: any) => {
                console.error('Error creating user', error);
                this.router.navigate(['/']);
            }
        });
    }

    updateFormWithUser(user: User){
        this.form = this.formBuilder.group({
            user_name: [user.user_name, Validators.required],
            first_name: [user.first_name, Validators.required],
            last_name: [user.last_name, Validators.required],
            email: [user.email, Validators.compose([Validators.required, Validators.email])],
            user_status: [user.user_status, Validators.required],
            department: [user.department],
          });
    }

    onSubmit(){
        if (this.form.valid) {
            const newUser: User = this.form.value;
            this.userService.createUser(newUser).subscribe({
                next: (data: any) => {
                    console.error('Data', data);
                    this.router.navigate(['/']);
                    this.form.reset();
                },
                error: (error: any) => {
                    console.error('Error creating user', error);
                    this.router.navigate(['/']);
                }
            });
        }
    }
}
  