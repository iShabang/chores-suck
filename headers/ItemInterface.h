/*############################################################################
 * File: headers/ItemInterface.h
 * Author: Shannon Montoya-Curtin
 * Date: 8/29/2019
############################################################################*/

#ifndef ITEM_INTERFACE_H
#define ITEM_INTERFACE_H

namespace item {

class ItemInterface {

    public:
        virtual ~ItemInterface() = 0;
        virtual void addItem() = 0;
        virtual void removeItem() = 0;

};

} // namespace item
#endif //ITEM_INTERFACE_H
