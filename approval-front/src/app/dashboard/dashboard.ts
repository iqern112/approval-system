import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { ApiService } from '../services/api';
import { FormsModule } from '@angular/forms';

export interface ApprovalResponse {
  id: number;
  title: string;
  admin_reason: string;
  status: string;
  username: string;
  created_at: string;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard implements OnInit {
  userRole = signal<string | null>(null);
  requests = signal<ApprovalResponse[]>([]);
  selectedIds = signal<Set<number>>(new Set());

  isModalOpen = signal(false);
  newRequestTitle = signal('');
  isSubmitting = signal(false);

  isActionModalOpen = signal(false);
  actionType = signal<'approved' | 'rejected' | null>(null);
  actionReason = signal('');
  pendingActionIds = signal<number[]>([]);

  statusLabel: Record<string, string> = {
    pending: 'รอนุมัติ',
    approved: 'อนุมัติ',
    rejected: 'ไม่อนุมัติ',
  };

  private router = inject(Router);
  private apiService = inject(ApiService);

  ngOnInit() {
    const role = localStorage.getItem('role');
    this.userRole.set(role);

    if (!role) {
      this.router.navigate(['/login']);
      return;
    }

    if (role === 'user') this.fetchMyRequests();
    else if (role === 'admin') this.fetchAllRequests();
  }

  fetchAllRequests() {
    this.apiService.getAllRequests().subscribe({
      next: (data: ApprovalResponse[]) => {
        this.requests.set(data);
        this.selectedIds.set(new Set());
      },
      error: (err) => {
        if (err.status === 401) this.onLogout();
      },
    });
  }

  fetchMyRequests() {
    this.apiService.getMyRequests().subscribe({
      next: (data: ApprovalResponse[]) => {
        this.requests.set(data);
      },
      error: (err) => {
        if (err.status === 401) this.onLogout();
      },
    });
  }

  // --- Checkbox ---

  toggleSelect(id: number) {
    const current = new Set(this.selectedIds());
    current.has(id) ? current.delete(id) : current.add(id);
    this.selectedIds.set(current);
  }

  toggleAll(event: Event) {
    const checked = (event.target as HTMLInputElement).checked;
    this.selectedIds.set(checked ? new Set(this.requests().map((r) => r.id)) : new Set());
  }

  isAllSelected(): boolean {
    return this.requests().length > 0 && this.selectedIds().size === this.requests().length;
  }

  isIndeterminate(): boolean {
    return this.selectedIds().size > 0 && this.selectedIds().size < this.requests().length;
  }

  // --- Modal ---

  openModal() {
    this.isModalOpen.set(true);
  }

  closeModal() {
    this.isModalOpen.set(false);
    this.newRequestTitle.set('');
  }

  submitRequest() {
    if (!this.newRequestTitle().trim()) {
      alert('กรุณากรอกรายละเอียดคำขอ');
      return;
    }

    this.isSubmitting.set(true);

    this.apiService.createRequest(this.newRequestTitle()).subscribe({
      next: () => {
        alert('ส่งคำขอสำเร็จแล้ว!');
        this.closeModal();
        this.fetchMyRequests();
        this.isSubmitting.set(false);
      },
      error: (err) => {
        console.error(err);
        alert('เกิดข้อผิดพลาดในการส่งคำขอ');
        this.isSubmitting.set(false);
      },
    });
  }

  openActionModal(type: 'approved' | 'rejected') {
    const ids = Array.from(this.selectedIds());
    if (ids.length === 0) return;

    this.pendingActionIds.set(ids);
    this.actionType.set(type);
    this.actionReason.set('');
    this.isActionModalOpen.set(true);
  }

  closeActionModal() {
  this.isActionModalOpen.set(false);
  this.actionType.set(null);
  this.actionReason.set('');
  this.pendingActionIds.set([]);
}
  confirmAction() {
  if (!this.actionReason().trim()) {
    alert('กรุณากรอกเหตุผลก่อนยืนยัน');
    return;
  }

  const ids = this.pendingActionIds();
  const newStatus = this.actionType()!;
  const reason = this.actionReason();


  this.apiService.approveRequests(ids, newStatus, reason).subscribe({
    next: (res) => {
      this.requests.set(
        this.requests().map((r) =>
          ids.includes(r.id) && r.status === 'pending'
            ? { ...r, status: newStatus, admin_reason: reason }
            : r
        )
      );
      this.selectedIds.set(new Set());
      this.closeActionModal();
      alert(res.message);
    },
    error: (err) => {
      console.error(err);
      alert('เกิดข้อผิดพลาด: ' + (err.error?.error || 'ไม่สามารถดำเนินการได้'));
    },
  });
}
  bulkAction(newStatus: 'approved' | 'rejected') {
    const ids = Array.from(this.selectedIds());
    if (ids.length === 0) return;

    this.requests.set(
      this.requests().map((r) =>

        ids.includes(r.id) && r.status === 'pending' ? { ...r, status: newStatus } : r,
      ),
    );
    this.selectedIds.set(new Set());
  }

  onLogout() {
    localStorage.clear();
    this.router.navigate(['/login']);
  }
}
