Meteor.startup(function() {
    // code to run on server at startup
    Meteor.publish('statements', function() { return Statements.find(); });
    Meteor.publish('predictions', function() { return Predictions.find(); });
});

// Allows people to see user information about others than themselves
Meteor.publish("allUsers", function () {
    return Meteor.users.find({}, {
        fields: {
            '_id': 1,
            'profile': 1,
            'username': 1
    }});
});

Votes.deny({
    "insert": function(userId, doc) {
        // TODO: This method for deleting existing votes is nasty,
        // move that logic to votebutton controller.
        existingVote = Votes.findOne({createdBy: userId, post: doc.post, type: "UpDown"});

        if(existingVote === undefined) {
            return false;
        } else {
            // Remove existing vote
            Votes.remove(existingVote._id);

            // Overwrite if vote with differing value already exists
            if(existingVote.value !== doc.value) {
                return false;
            }

            // If previous vote was identical in value, don't insert it.
            return true;
        }
    }
});
