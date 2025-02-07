<template>
  <div class="login-container">
    <h2>Login</h2>
    <form @submit.prevent="loginUser">
      <div class="mb-3">
        <label for="username" class="form-label d-block">Username</label>
        <input 
          type="text" 
          id="username" 
          class="form-control-md" 
          v-model="username" 
          placeholder="Enter your username" 
          required
        />
      </div>
      <button type="submit" class="btn btn-primary">Login</button>
      <p v-if="msg" class="text-danger">{{ msg }}</p>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      msg: "",
      securityKey: null,
      userId: null
    };
  },
  methods: {
    async loginUser() {
      this.msg = "";
      try {
        // Login to get the security key and user ID
        let response = await this.$axios.post("/session", { name: this.username });

        // Store security key and user ID
        this.securityKey = response.data.apiKey;
        this.userId = response.data.userId;

        // Fetch username using the user ID and security key
        let userResponse = await this.$axios.get(`/users/${this.userId}/username`, {
          headers: { Authorization: `Bearer ${this.securityKey}` }
        });

        this.username = userResponse.data.username; // Display the retrieved username
        this.msg = "Logged in successfully";
      } catch (e) {
        this.msg = "Login failed: " + (e.response?.data?.error || e.message);
      }
    }
  }
};
</script>


  

<style>
</style>
