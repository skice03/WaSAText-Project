<template>
	<div v-if="isAuthorized">
	  <h2>Welcome to WaSAText</h2>
	  <p v-if="!loading && chats.length === 0">No chats available.</p>
	  <button class="btn btn-primary mb-3" @click="toggleCreateChat">
		{{ showCreateChat ? "Cancel" : "Create New Chat" }}
	  </button>
	  
	  <div v-if="showCreateChat" class="create-chat-form">
		<h4>Create a New Chat</h4>
		<form @submit.prevent="createChat">
		  <div class="mb-3">
			<label for="chatName" class="form-label">Chat Name</label>
			<input type="text" id="chatName" v-model="newChatName" class="form-control" required />
		  </div>
		  <button type="submit" class="btn btn-success">Create Chat</button>
		</form>
	  </div>
  
	  <div v-if="loading" class="loading-spinner">Loading chats...</div>
	  <ul v-if="!loading && chats.length > 0" class="chat-list">
		<li v-for="chat in chats" :key="chat.id" @click="openChat(chat.id)" class="chat-item">
		  {{ chat.name }} <small>(Type: {{ chat.type }})</small>
		</li>
	  </ul>
  
	  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
  
	<div v-else>
	  <p class="text-danger">You are not logged in. Please log in to access the home page.</p>
	</div>
  </template>
  
  <script>
  export default {
	data() {
	  return {
		chats: [],
		newChatName: "",
		showCreateChat: false,
		loading: true,
		errormsg: "",
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
  
		// Verify the token by making a GET request
		await this.$axios.get(`/users/${userId}/username`, {
		  headers: { Authorization: `Bearer ${securityKey}` }
		});
  
		this.isAuthorized = true;
		await this.fetchChats();
	  } catch (e) {
		this.isAuthorized = false;
		console.error("Authorization failed:", e.response?.data?.error || e.message);
	  }
	},
	methods: {
	  async fetchChats() {
		try {
		  this.loading = true;
		  const response = await this.$axios.get("/chats", {
			headers: { Authorization: `Bearer ${localStorage.getItem("securityKey")}` }
		  });
		  this.chats = response.data.userChats || [];
		} catch (e) {
		  this.errormsg = "Failed to load chats: " + (e.response?.data?.error || e.message);
		} finally {
		  this.loading = false;
		}
	  },
	  toggleCreateChat() {
		this.showCreateChat = !this.showCreateChat;
		this.newChatName = "";
	  },
	  async createChat() {
		try {
		  await this.$axios.post(
			"/chats",
			{ name: this.newChatName },
			{ headers: { Authorization: `Bearer ${localStorage.getItem("securityKey")}` } }
		  );
		  this.newChatName = "";
		  this.showCreateChat = false;
		  await this.fetchChats();
		} catch (e) {
		  this.errormsg = "Failed to create chat: " + (e.response?.data?.error || e.message);
		}
	  },
	  openChat(chatId) {
		this.$router.push(`/chats/${chatId}`);
	  }
	}
  };
  </script>
  
  <style>

  .chat-list {
	list-style: none;
	padding: 0;
  }
  
  .chat-item {
	padding: 10px;
	border-bottom: 1px solid #ccc;
	cursor: pointer;
  }
  
  .chat-item:hover {
	background-color: #f0f0f0;
  }
  
  .loading-spinner {
	text-align: center;
	font-weight: bold;
  }
  </style>
  