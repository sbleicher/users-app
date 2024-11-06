import { Component, OnInit } from '@angular/core';
import { NgFor, NgIf } from '@angular/common';
import {MatTableModule} from '@angular/material/table';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {User, UserStatus, userStatusToString} from '../user';
import { Router } from '@angular/router';
import { UserService } from '../services/user.service';

@Component({
  standalone: true,
  selector: 'users-list',
  templateUrl: './user-list.component.html',
  imports: [NgFor, NgIf, MatTableModule, MatIconModule, MatButtonModule],
})
export class UserListComponent implements OnInit {
    displayedColumns: string[] = ['user_id', 'user_name', 'first_name', 'last_name', 'email', 'user_status', 'department', 'actions'];
    users: User[] = [];

  constructor(private router: Router, private userService: UserService) {}

  ngOnInit() {
    this.userService.getUsers().subscribe({
      next: (req: any) => {
        this.users = req.data ? req.data : []
      },
      error: (error: any) => {
        console.error('Error getting users', error);
      }
    });
  }

  navigateToCreatePage() {
    this.router.navigate(['/create']);
  }

  navigateToEditPage(user_id: number) {
    this.router.navigate(['/edit'], { queryParams: { user_id } });
  }

  deleteUser(user_id: number) {
    this.userService.deleteUser(user_id).subscribe({
      next: () => {
        this.reloadPage();
      },
      error: (error: any) => {
        console.error('Error deleting user', error);
        window.location.reload();
      }
    });
  }

  statusToString(user_status: UserStatus): string {
    return userStatusToString(user_status);
  }

  reloadPage() {
    window.location.reload();
  }
}

