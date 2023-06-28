package routes

import (
	"github.com/francischacko/ecommerce/controllers"
	middleware "github.com/francischacko/ecommerce/middlewares"

	"github.com/gin-gonic/gin"
)

func UserInfo(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/signup", controllers.Signup)
		user.POST("/login", controllers.Login)

		//otp verification
		user.POST("/sendOTP", controllers.Verify)
		user.POST("/getOTP", controllers.CheckOtp)

		//Product listing
		user.GET("/listproducts", middleware.UserAuth, controllers.ListProducts)
		//cart management of user
		user.POST("/addtocart", middleware.UserAuth, controllers.AddToCart)
		user.GET("/listcart", middleware.UserAuth, controllers.ListCart)
		user.PATCH("/quantityupdation", middleware.UserAuth, controllers.QuantityUpdation)
		user.DELETE("/cartitemdeletion", middleware.UserAuth, controllers.CartItemDeletion)
		//user profile management
		user.POST("/addaddress", middleware.UserAuth, controllers.AddUserAddress)
		user.PATCH("/editaddress", middleware.UserAuth, controllers.EditUserAddress)
		user.GET("/showuseraddress", middleware.UserAuth, controllers.ShowUserAddress)
		user.GET("/chooseaddress", middleware.UserAuth, controllers.ChooseAddress)
		user.PATCH("/changepassword", middleware.UserAuth, controllers.ChangePassword)
		//user order management
		user.GET("/listordersuser", middleware.UserAuth, controllers.ListOrders)
		user.POST("/placeorder", middleware.UserAuth, controllers.PlaceOrder)
		user.POST("/cancelorderuser", middleware.UserAuth, controllers.CancelOrder)
		user.PUT("/return", middleware.UserAuth, controllers.ReturnOrder)

		// coupen management
		user.POST("/coupen", middleware.UserAuth, controllers.RedeemCoupen)
		user.POST("/coupenvalidation", middleware.UserAuth, controllers.CoupenValidation)

		//Payment
		user.GET("/payment-success", middleware.UserAuth, controllers.RazorpaySuccess)
		user.GET("/success", middleware.UserAuth, controllers.Success)
		user.GET("/razorpay", middleware.UserAuth, controllers.RazorPay)

	}
}

func AdminInfo(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.POST("/signup", controllers.SignupAdmin)
		admin.POST("/login", controllers.LoginAdmin)

		//user management
		admin.GET("/getall", middleware.AdminAuth, controllers.GetAllUsers)
		admin.PATCH("/block/:id", middleware.AdminAuth, controllers.BlockUser)
		admin.PATCH("/unblock/:id", middleware.AdminAuth, controllers.UnblockUser)
		//category management
		admin.POST("/addcategory", middleware.AdminAuth, controllers.AddCategory)
		admin.GET("/listcategories", middleware.AdminAuth, controllers.ListCategories)
		admin.PATCH("/editcategory/:id", middleware.AdminAuth, controllers.EditCategory)
		admin.DELETE("/deletecategory", middleware.AdminAuth, controllers.DeleteCategory)
		//admin side management of products
		admin.POST("/addproduct", middleware.AdminAuth, controllers.AddProduct)
		admin.GET("/listproducts", middleware.AdminAuth, controllers.ListProducts)
		admin.PATCH("/editproduct", middleware.AdminAuth, controllers.EditProduct)
		admin.DELETE("/deleteproduct", middleware.AdminAuth, controllers.DeleteProduct)
		//admin side management of orders
		admin.GET("/listorders", middleware.AdminAuth, controllers.ListOrders)
		admin.POST("/cancelorder", middleware.AdminAuth, controllers.CancelOrder)
		// banner mangement
		admin.POST("/addbanner", middleware.AdminAuth, controllers.AddBanner)
		//coupen management
		admin.POST("/addcoupen", middleware.AdminAuth, controllers.AddCoupen)
	}
}
