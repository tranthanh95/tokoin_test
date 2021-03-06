package services

import (
	"fmt"

	"tokoin-challenge/src/entity"
	"tokoin-challenge/src/repositories"
)

type ITicketService interface {
	List(key, value string) (*entity.TicketsResponse, error)
}

type TicketService struct {
	orgRepo    repositories.IOrgRepository
	ticketRepo repositories.ITicketRepository
	userRepo   repositories.IUserRepository
}

func NewTicketService(orgRepo repositories.IOrgRepository, ticketRepo repositories.ITicketRepository,
	userRepo repositories.IUserRepository) ITicketService {
	return &TicketService{
		orgRepo:    orgRepo,
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

func (s *TicketService) List(key, value string) (*entity.TicketsResponse, error) {
	tickets, err := s.ticketRepo.List(key, value)
	if err != nil {
		return nil, err
	}
	results := entity.TicketsResponse{}
	for _, ticket := range *tickets {
		var rs entity.TicketResponse
		rs.Ticket = *ticket

		// Get organization of ticket
		org, err := s.orgRepo.Retrieve(ticket.OrganizationID)
		if err != nil {
			fmt.Printf("Cannot get organization of ticket %s. Error: %s\n", ticket.ID, err)
		}
		rs.OrganizationName = org.Name

		// Get assignee for ticket
		assignee, err := s.userRepo.Retrieve(ticket.AssigneeID)
		if err != nil {
			fmt.Printf("Cannot get assignee of ticket %s. Error: %s\n", ticket.ID, err)
		}
		if assignee != nil {
			rs.AssigneeName = assignee.Name
		}

		// Get submitter for ticket
		submitter, err := s.userRepo.Retrieve(ticket.SubmitterID)
		if err != nil {
			fmt.Printf("Cannot get submitter of ticket %s. Error: %s\n", ticket.ID, err)
		}
		if submitter != nil {
			rs.SubmitterName = submitter.Name
		}

		results = append(results, &rs)
	}

	return &results, nil
}
