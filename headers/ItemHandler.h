/*############################################################################
 * File: headers/ItemHandler.h
 * Author: Shannon Montoya-Curtin
 * Date: 8/29/2019
############################################################################*/

#ifndef ITEM_HANDLER_H
#define ITEM_HANDLER_H

#include "headers/ItemInterface.h"

namespace item {

class ItemHandler : public ItemInterface
{
    public:
        ItemHandler();

        void addItem();
        void removeItem();

};

} // namespace item
#endif // ITEM_HANDLER_H
