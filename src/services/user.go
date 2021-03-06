package services

import (
	"fmt"
	"strconv"

	"tokoin-challenge/src/entity"
	"tokoin-challenge/src/repositories"
)

type IUserService interface {
	List(key, value string) (*entity.UsersResponse, error)
}

type UserService struct {
	orgRepo    repositories.IOrgRepository
	ticketRepo repositories.ITicketRepository
	userRepo   repositories.IUserRepository
}

func NewUserService(orgRepo repositories.IOrgRepository, ticketRepo repositories.ITicketRepository,
	userRepo repositories.IUserRepository) IUserService {
	return &UserService{
		orgRepo:    orgRepo,
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

func (s *UserService) List(key, value string) (*entity.UsersResponse, error) {
	users, err := s.userRepo.List(key, value)
	if err != nil {
		return nil, err
	}
	results := entity.UsersResponse{}
	for _, u := range *users {
		var rs entity.UserResponse

		rs.User = *u
		strUID := strconv.Itoa(u.ID)

		// Get assignee tickets tickets for user
		rs.AssigneeTicketSubjects, _ = s.ticketRepo.ListSubjects("assignee_id", strUID)

		// Get submitted tickets tickets for user
		rs.SubmittedTicketSubjects, _ = s.ticketRepo.ListSubjects("submitter_id", strUID)

		// Get organization of user
		org, err := s.orgRepo.Retrieve(u.OrganizationID)
		if err != nil {
			fmt.Printf("Cannot get organization of user %d. Error: %s\n", u.ID, err)
		}
		rs.OrganizationName = org.Name

		results = append(results, &rs)
	}

	return &results, nil
}
