package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["airtime_payment_service/controllers:CallbackController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:CallbackController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/process`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "GetBundles",
            Router: `/bundles/:networkId/:destinationPhoneNumber`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "BuyAirtime",
            Router: `/buy-airtime`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:RequestController"],
        beego.ControllerComments{
            Method: "BuyDataBundle",
            Router: `/buy-bundle`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"] = append(beego.GlobalControllerRouter["airtime_payment_service/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
