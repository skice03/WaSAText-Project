<template>
    <div v-if="isAuthorized">
      <h2>My Profile</h2>
      <form @submit.prevent="updateUsername">
        <div class="mb-3">
          <label for="newUsername" class="form-label d-block">New Username</label>
          <input 
            type="text" 
            id="newUsername" 
            class="form-control-md" 
            v-model="newUsername" 
            placeholder="Enter your new username" 
            required
          />
        </div>
        <button type="submit" class="btn btn-success">Update Username</button>
        <p v-if="msg" :class="{ 'text-success': success, 'text-danger': !success }">{{ msg }}</p>
      </form>
    </div>
    <div v-else>
      <p class="text-danger">You are not logged in. Please log in to access your profile.</p>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        newUsername: "",
        msg: "",
        success: false,
        isAuthorized: false
      };
    },
    async created() {
      try {
        const userId = localStorage.getItem("userId");
        const securityKey = localStorage.getItem("securityKey");
  
        if (!userId || !securityKey) {
          throw new Error("Authorization data is missing.");
        }
  
        // Veifying the token
        await this.$axios.get(`/users/${userId}/username`, {
          headers: { Authorization: `Bearer ${securityKey}` }
        });
  
        this.isAuthorized = true;
      } catch (e) {
        this.isAuthorized = false;
        console.error("Authorization failed:", e.response?.data?.error || e.message);
      }
    },
    methods: {
      async updateUsername() {
        this.msg = "";
        this.success = false;
        try {
          const userId = localStorage.getItem("userId");
          const securityKey = localStorage.getItem("securityKey");
  
          let response = await this.$axios.put(
            `/users/${userId}/username`,
            { username: this.newUsername },
            {
              headers: { Authorization: `Bearer ${securityKey}` }
            }
          );
  
          if (response.status === 204) {
            this.msg = "Username updated successfully.";
            this.success = true;
          }
        } catch (e) {
          this.msg = `Failed to update username: ${e.response?.data?.error || e.message}`;
        }
      }
    }
  };
  </script>
  
  <style>
  </style>
  