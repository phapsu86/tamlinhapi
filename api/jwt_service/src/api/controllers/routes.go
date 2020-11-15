package controllers

import "github.com/victorsteven/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/logout", middlewares.SetMiddlewareJSON(s.Logout)).Methods("POST")
	s.Router.HandleFunc("/forgot_pass", middlewares.SetMiddlewareJSON(s.ForgotPass)).Methods("POST")
	s.Router.HandleFunc("/profile", middlewares.SetMiddlewareJSON(s.GetProfileUser)).Methods("POST")
	s.Router.HandleFunc("/change_pass", middlewares.SetMiddlewareJSON(s.ChangePassword)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/user_rsa", middlewares.SetMiddlewareJSON(s.CreateUserRSA)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/update_users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("POST")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Posts routes
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/religion_post/{id}/{p}", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/religion_post_suggest/{id}/{post_id}/{p}", middlewares.SetMiddlewareJSON(s.GetSuggestPosts)).Methods("GET")
	s.Router.HandleFunc("/post_comment", middlewares.SetMiddlewareJSON(s.CreatePostComment)).Methods("POST")
	s.Router.HandleFunc("/post_comment/{id}/{p}", middlewares.SetMiddlewareJSON(s.GetPostCommentByID)).Methods("GET")

	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")
	// ReligionPOSTLsfj   like folow share Join cho post
	s.Router.HandleFunc("/religion_post_lsfcs", middlewares.SetMiddlewareJSON(s.CreatePostLsfc)).Methods("POST")
	s.Router.HandleFunc("/religion_post_lsfcs/{id}", middlewares.SetMiddlewareJSON(s.GetTotalPostLsfjByID)).Methods("GET")
	s.Router.HandleFunc("/get_post_user_lsfjs/{id}", middlewares.SetMiddlewareJSON(s.GetPostLSFJByUser)).Methods("GET") //Lấy thông tin user lsfj cho 1 post
	s.Router.HandleFunc("/get_all_post_user_lsc", middlewares.SetMiddlewareJSON(s.GetAllPostLSFJByUser)).Methods("GET") //Lấy thông tin user lsfj cho 1 post

	//ReligionListRouter
	s.Router.HandleFunc("/religions", middlewares.SetMiddlewareJSON(s.CreateReligionList)).Methods("POST")
	s.Router.HandleFunc("/religions", middlewares.SetMiddlewareJSON(s.GetReligionLists)).Methods("GET")
	s.Router.HandleFunc("/religions/{id}", middlewares.SetMiddlewareJSON(s.GetReligionList)).Methods("GET")
	s.Router.HandleFunc("/religions/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateReligionList))).Methods("PUT")
	s.Router.HandleFunc("/religions/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteReligionList)).Methods("DELETE")

	//ReligionItem
	s.Router.HandleFunc("/religion_items", middlewares.SetMiddlewareJSON(s.GetReligionItems)).Methods("GET")
	s.Router.HandleFunc("/religion_items/{id}/{page}", middlewares.SetMiddlewareJSON(s.GetReligionItemByReligionID)).Methods("GET")
	s.Router.HandleFunc("/religion_item_details/{id}", middlewares.SetMiddlewareJSON(s.GetReligionItemDetails)).Methods("GET")
	s.Router.HandleFunc("/religion_item_search", middlewares.SetMiddlewareJSON(s.SearchItemByNames)).Methods("POST")
	s.Router.HandleFunc("/religion_item_follows", middlewares.SetMiddlewareJSON(s.CreatItemFollow)).Methods("POST")
	s.Router.HandleFunc("/get_religion_item_follow", middlewares.SetMiddlewareJSON(s.GetReligionItemFollow)).Methods("GET")
	s.Router.HandleFunc("/get_info_item_follow_by_id/{id}", middlewares.SetMiddlewareJSON(s.GetInfoItemFollowByID)).Methods("GET")
	s.Router.HandleFunc("/check_user_of_item/{id}", middlewares.SetMiddlewareJSON(s.CheckUserOfItem)).Methods("GET")

	//ReligionItemEvent by religion_id

	s.Router.HandleFunc("/religion_item_events/{id}/{type}/{page}", middlewares.SetMiddlewareJSON(s.GetReligionItemEventByID)).Methods("GET")
	s.Router.HandleFunc("/religion_events/{id}/{page}", middlewares.SetMiddlewareJSON(s.GetReligionEventByID)).Methods("GET")
	s.Router.HandleFunc("/religion_event_details/{id}", middlewares.SetMiddlewareJSON(s.GetReligionEventDetails)).Methods("GET")

	s.Router.HandleFunc("/religion_event_search_by_name", middlewares.SetMiddlewareJSON(s.SearchEventByNames)).Methods("POST")
	s.Router.HandleFunc("/religion_event_search_by_date", middlewares.SetMiddlewareJSON(s.SearchEventByDate)).Methods("POST")
	// ReligionEventLsfj   like folow share Join cho event
	s.Router.HandleFunc("/religion_event_lsfjs", middlewares.SetMiddlewareJSON(s.CreateReligionEventLsfj)).Methods("POST")
	s.Router.HandleFunc("/religion_event_lsfjs/{id}", middlewares.SetMiddlewareJSON(s.GetTotalEventLsfjByID)).Methods("GET")
	s.Router.HandleFunc("/religion_event_lsfjs_follow/{page}", middlewares.SetMiddlewareJSON(s.GetEventLSFJFollow)).Methods("GET")
	s.Router.HandleFunc("/get_event_user_lsfjs/{id}", middlewares.SetMiddlewareJSON(s.GetInfoEventLSFJByID)).Methods("GET") //Lấy thông tin user lsfj cho 1 sự kiện

	// Cong duc
	s.Router.HandleFunc("/religion_merit_list/{id}/{type}/{page}", middlewares.SetMiddlewareJSON(s.GetMeritList)).Methods("GET")
	s.Router.HandleFunc("/religion_merits", middlewares.SetMiddlewareJSON(s.CreateMerit)).Methods("POST")                     //Tao cong duc
	s.Router.HandleFunc("/religion_merit_address", middlewares.SetMiddlewareJSON(s.CreateMeritAddress)).Methods("POST")       //Tao cong duc
	s.Router.HandleFunc("/religion_merit_event_search", middlewares.SetMiddlewareJSON(s.SearchEventForMerit)).Methods("POST") //Tìm kiếm sự kiện cho công đức
	s.Router.HandleFunc("/search_offering_item", middlewares.SetMiddlewareJSON(s.SearchOfferItemByNames)).Methods("POST")     //Tìm kiếm sự kiện cho công đức

	//Danh sach vat pham cho sukien hoaj len hoi
	s.Router.HandleFunc("/religion_offering_list/{id}/{type}/{page}", middlewares.SetMiddlewareJSON(s.GetOfferingItemSells)).Methods("GET")
	s.Router.HandleFunc("/religion_merit_user_info", middlewares.SetMiddlewareJSON(s.GetMeritUserInfo)).Methods("POST")
	// Lây point của user
	s.Router.HandleFunc("/points", middlewares.SetMiddlewareJSON(s.GetPoint)).Methods("POST")

	//UPLOAD FILE
	s.Router.HandleFunc("/add_files", middlewares.SetMiddlewareJSON(s.UploadFile)).Methods("POST")

	//Remind User
	s.Router.HandleFunc("/add_event_remind", middlewares.SetMiddlewareJSON(s.CreateRemind)).Methods("POST")
	s.Router.HandleFunc("/list_remind", middlewares.SetMiddlewareJSON(s.GetListUserRemind)).Methods("GET")

	//OTP

	s.Router.HandleFunc("/send_otp", middlewares.SetMiddlewareJSON(s.SendOTP)).Methods("POST")
	s.Router.HandleFunc("/get_otp/{phone}", middlewares.SetMiddlewareJSON(s.GetOTP)).Methods("GET")
	//Get Info Banks
	s.Router.HandleFunc("/get_info_bank/{id}/{type}", middlewares.SetMiddlewareJSON(s.GetInfoBank)).Methods("GET")
	s.Router.HandleFunc("/get_info_bank_address", middlewares.SetMiddlewareJSON(s.GetInfoBankAddress)).Methods("GET")

}
