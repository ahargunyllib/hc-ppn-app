const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export interface PaginationMeta {
  total_data: number;
  total_page: number;
  page: number;
  limit: number;
}

export interface FeedbacksResponse {
  feedbacks: Array<{
    id: string;
    sessionId: string;
    phoneNumber: string;
    rating: number;
    comment?: string;
    createdAt: string;
  }>;
  meta: {
    pagination: PaginationMeta;
  };
}

export interface PhoneNumbersResponse {
  users: Array<{
    id: string;
    phoneNumber: string;
    label: string;
    assignedTo?: string;
    notes?: string;
    createdAt: string;
    updatedAt: string;
  }>;
  meta: {
    pagination: PaginationMeta;
  };
}

export const apiClient = {
  async getFeedbacks(params: { page?: number; limit?: number; phoneNumber?: string; minRating?: number; maxRating?: number }): Promise<FeedbacksResponse> {
    const queryParams = new URLSearchParams();

    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());
    if (params.phoneNumber) queryParams.append("phoneNumber", params.phoneNumber);
    if (params.minRating) queryParams.append("minRating", params.minRating.toString());
    if (params.maxRating) queryParams.append("maxRating", params.maxRating.toString());

    const url = `${API_BASE_URL}/api/v1/feedbacks?${queryParams.toString()}`;

    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch feedbacks: ${response.statusText}`);
    }

    const data = await response.json();
    return data.data;
  },

  async getPhoneNumbers(params: { page?: number; limit?: number; search?: string; assignedTo?: string }): Promise<PhoneNumbersResponse> {
    const queryParams = new URLSearchParams();

    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());
    if (params.search) queryParams.append("search", params.search);
    if (params.assignedTo) queryParams.append("assignedTo", params.assignedTo);

    const url = `${API_BASE_URL}/api/v1/users?${queryParams.toString()}`;

    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch phone numbers: ${response.statusText}`);
    }

    const data = await response.json();
    return data.data;
  },
};
