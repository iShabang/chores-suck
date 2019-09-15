#include <string>
#include <vector>
#include <iostream>
#include <thread>
#include <chrono>

enum TimeType{
    NONE,
    DAY,
    WEEK,
    MONTH,
    YEAR
};

struct ItemStructure{
    std::string name;
    unsigned int startDay;
    unsigned int endDay;
};

void printItems(const std::vector<ItemStructure> &items, unsigned int currentDay);
void checkItemsExpired(const std::vector<ItemStructure> &items, unsigned int currentDay);
int printMenu();
TimeType findTimeType(std::string type);
unsigned int convertTime(unsigned int amount, TimeType timeType);

int main(){
    std::vector<ItemStructure> items;
    int option = 1;
    std::string itemString("");
    unsigned int days(0);
    unsigned int currentDay(0); 
    TimeType timeType;

    /*
     * Setup for time arithmetic
    std::chrono::system_clock::time_point now = std::chrono::system_clock::now();
    std::time_t startTime = std::chrono::system_clock::to_time_t(now);
    std::time_t currentTime = std::chrono::system_clock::to_time_t(now);
    std::time_t newTime = currentTime - startTime;
    */

    while(option > 0){
        checkItemsExpired(items,currentDay);
        option = printMenu();
        ItemStructure temp;
        temp.name = "";
        temp.startDay = 0;
        temp.endDay = 0;
        switch(option){
            case 1:
                std::cout << "name:" << std::endl;
                std::cin >> itemString;
                temp.name = itemString;
                std::cout << "Enter time type: " << std::endl;
                std::cin >> itemString;
                std::cout << "Enter time amount:" << std::endl;
                std::cin >> days;
                timeType = findTimeType(itemString);
                temp.startDay = currentDay;
                days = convertTime(days, timeType);
                temp.endDay = temp.startDay + days;
                items.push_back(temp);
                break;
            case 2:
                currentDay++;
                printItems(items, currentDay);
                break;
            default:
                break;
        }
    }
}

void printItems(const std::vector<ItemStructure> &items, unsigned int currentDay){
    if (!items.size()){
        return;
    }
    std::vector<ItemStructure>::const_iterator iter = items.begin();  
    for (; iter != items.end(); iter++){
        std::cout << iter->name << "   " << iter->endDay  - currentDay << " days remaining" << std::endl;
    }
}

void checkItemsExpired(const std::vector<ItemStructure> &items, unsigned int currentDay){
    std::vector<ItemStructure>::const_iterator iter;
    for (iter = items.begin(); iter != items.end(); iter++){
        if (iter->endDay - currentDay < 1){
            std::cout << iter->name << " is expired! Calling subroutine!\n";
        }
    }
}

int printMenu(){
    int option(0);
    std::cout << "Menu \n";
    std::cout << "1: Add Item\n";
    std::cout << "2: Advance Day\n";
    std::cout << "0: Exit\n";
    std::cout << "Enter option: ";
    std::cin >> option;
    return option;
}

TimeType findTimeType(std::string type){
    if (type == "d" || type == "D")
        return DAY;
    if (type == "w" || type == "W")
        return WEEK;
    return NONE;
}

unsigned int convertTime(unsigned int amount, TimeType timeType){
    unsigned int days(0);
    switch (timeType){
        case DAY:
            days = amount;
            break;
        case WEEK:
            days = amount * 7;
            break;
        default:
            break;
    }
    return days;
}
