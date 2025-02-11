<template>
	<div v-if="isAuthorized">
	  <h2>Welcome to WaSAText</h2>
	  <p v-if="!loading && chats.length === 0">No chats available.</p>
  
	  <button class="btn btn-primary mb-3" @click="toggleCreateChat">
		{{ showCreateChat ? "Cancel" : "Create New Chat" }}
	  </button>
  
	  <div v-if="showCreateChat" class="create-chat-form">
		<h4>Create a new group chat</h4>
		<form @submit.prevent="createGroupChat">
		  <div class="mb-3">
			<label for="chatName" class="form-label">Chat name</label>
			<input type="text" id="chatName" v-model="newChatName" class="form-control" required />
		  </div>
  
		  <div class="mb-3">
			<label class="form-label">Add users to the conversation</label>
			<ul v-if="availableUsers.length > 0" class="user-list">
			  <li v-for="user in availableUsers" :key="user.id">
				<label>
				  <input type="checkbox" :value="user.id" v-model="selectedMembers" />
				  {{ user.name }}
				</label>
			  </li>
			</ul>
			<p v-if="!loading && availableUsers.length === 0" class="text-danger">
			  No users available to add to this chat.
			</p>
		  </div>
  
		  <button type="submit" class="btn btn-success" :disabled="selectedMembers.length === 0 || availableUsers.length === 0">
			Create Group Chat
		  </button>
		</form>
	  </div>
  
	  <div v-if="loading" class="loading-spinner">Loading chats...</div>
	  <ul v-if="!loading && chats.length > 0" class="chat-list">
		<li v-for="chat in chats" :key="chat.id" @click="openChat(chat.id)" class="chat-item">
		  {{ getChatLabel(chat) }}
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
		availableUsers: [],
		selectedMembers: [],
		loading: true,
		errormsg: "",
		isAuthorized: false,
		currentUserId: null,
		currentUsername: "",
	  };
	},
	async created() {
	  try {
		const userId = localStorage.getItem("userId");
		const securityKey = localStorage.getItem("securityKey");
  
		if (!userId || !securityKey) {
		  throw new Error("Authorization data is missing.");
		}
  
		const userResponse = await this.$axios.get(`/users/${userId}/username`, {
		  headers: { Authorization: `Bearer ${securityKey}` },
		});
  
		this.isAuthorized = true;
		this.currentUserId = userId;
		this.currentUsername = userResponse.data.username;
  
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
			headers: { Authorization: `Bearer ${localStorage.getItem("securityKey")}` },
		  });
		  this.chats = response.data.userChats || [];
		} catch (e) {
		  this.errormsg = "Failed to load chats: " + (e.response?.data?.error || e.message);
		} finally {
		  this.loading = false;
		}
	  },
	  async fetchAvailableUsers() {
		try {
		  const response = await this.$axios.get("/users", {
			headers: { Authorization: `Bearer ${localStorage.getItem("securityKey")}` },
		  });
		  this.availableUsers = response.data.filter(user => user.id !== this.currentUserId);
		} catch (e) {
		  console.warn("No other users available.");
		  this.availableUsers = [];
		}
	  },
	  toggleCreateChat() {
		this.showCreateChat = !this.showCreateChat;
		this.newChatName = "";
		this.selectedMembers = [];
  
		if (this.showCreateChat) {
		  this.availableUsers = [];
		  this.fetchAvailableUsers();
		}
	  },
	  async createGroupChat() {
		try {
		  const securityKey = localStorage.getItem("securityKey");
		  const createChatResponse = await this.$axios.post(
			"/chats",
			{ name: this.newChatName, groupChat: true },
			{ headers: { Authorization: `Bearer ${securityKey}` } }
		  );
  
		  const chatId = createChatResponse.data.chatId;
  
		  await this.$axios.post(
			`/chats/${chatId}/add`,
			{ members: this.selectedMembers },
			{ headers: { Authorization: `Bearer ${securityKey}` } }
		  );
  
		  this.newChatName = "";
		  this.selectedMembers = [];
		  this.showCreateChat = false;
		  await this.fetchChats();
		  this.errormsg = "Group chat created and members added successfully!";
		} catch (e) {
		  this.errormsg = "Failed to create group chat: " + (e.response?.data?.error || e.message);
		}
	  },
	  openChat(chatId) {
		this.$router.push(`/chats/${chatId}`);
	  },
	  getChatLabel(chat) {
		const memberNames = chat.members.filter(member => member.id !== this.currentUserId).map(member => member.name);
		if (memberNames.length === 1) {
		  return `Chat with ${memberNames[0]}`;
		}
		return `Chat with ${memberNames.join(", ")}`;
	  },
	},
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
  
  .user-list {
	list-style: none;
	padding: 0;
  }
  
  .user-list li {
	margin-bottom: 10px;
  }
  
  .text-danger {
	color: red;
  }
  </style>
  