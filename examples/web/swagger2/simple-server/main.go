package main

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/swagger2"
)

type employeeController struct {
	at.RestController
	at.RequestMapping `value:"/employee"`
	at.Api `value:"Employee Management System" description:"Operations pertaining to employee in Employee Management System"`
}

// newEmployeeController is the constructor for orgController
// you may inject dependency to constructor
func newEmployeeController() *employeeController {
	return new(employeeController)
}

// GetAnEmployee
// at.GetMapping is an annotation to define request mapping for http method GET,
func (c *employeeController) Get(at struct {
	at.GetMapping `value:"/{id:int}"`
	at.ApiOperation `value:"Get an employee"`
	at.ApiParam `value:"Path variable employee ID" required:"true"`
	// TODO: Code is not injected ...
	at.ApiResponse200 `value:"Successfully get an employee"`
	at.ApiResponse404 `value:"The resource you were trying to reach is not found"`
}, id int) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)
	response.Code = at.ApiResponse200.Code
	response.Message = "Success"
	response.Data = Employee{
		Id: id,
		FirstName: "John",
		LastName: "Deng",
	}
	return
}

// ListEmployee
func (c *employeeController) ListEmployee(at struct{ at.GetMapping `value:"/"` }, factory factory.ConfigurableFactory) (response model.Response, err error) {
	response = new(model.BaseResponse)
	return
}

func init() {
	app.Register(newEmployeeController)
}

// Hiboot main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude,
			actuator.Profile,
			swagger2.Profile,
		).Run()
}
