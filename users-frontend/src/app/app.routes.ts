import { RouterModule, Routes } from '@angular/router';
import { UserCreateUpdateComponent } from './components/user-create-update.component';
import { UserListComponent } from './components/users-list.component';
import { NgModule } from '@angular/core';

export const routes: Routes = [
    { path: '', component: UserListComponent },
    { path: 'create', component: UserCreateUpdateComponent },
    { path: 'edit', component: UserCreateUpdateComponent },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
